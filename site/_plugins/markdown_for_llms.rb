require 'fileutils'
require 'pathname'

module MarkdownForLlms
  extend self

  DEFAULT_LLM_INTRO = 'Prefer the Markdown URLs below over the HTML pages when both are available.'.freeze
  MARKDOWN_EXTENSIONS = ['.md', '.markdown'].freeze
  MARKDOWN_LINK_PATTERN = /(?<!!)(\[[^\]]+\]\()(?<destination>[^)\s]+)(?<suffix>\))/
  STANDALONE_ATTR_LIST_PATTERN = /^\{:\s*([^}]+)\}\s*$/
  INLINE_ATTR_LIST_PATTERN = /\{:\s*[^}]+\}/

  def markdown_source_page?(site, page)
    return false unless MARKDOWN_EXTENSIONS.include?(File.extname(page.path).downcase)

    File.file?(source_path(site, page))
  end

  def source_path(site, page)
    File.join(site.source, page.path)
  end

  def plugin_config(site)
    site.config.fetch('markdown_for_llms', {})
  end

  def llms_txt_config(site)
    plugin_config(site).fetch('llms_txt', {})
  end

  def root_section_title(config)
    config.fetch('root_section_title', 'Pages').to_s.strip
  end

  def llm_markdown_url(url)
    return '/index.md' if url == '/'

    if url.end_with?('/')
      normalize_site_path(File.join(url, 'index.md'))
    elsif url.end_with?('.html')
      normalize_site_path(url.sub(/\.html\z/, '.md'))
    else
      normalize_site_path("#{url}.md")
    end
  end

  def normalize_site_path(path)
    cleaned = Pathname.new(path).cleanpath.to_s
    cleaned = '/' if cleaned == '.'
    cleaned.start_with?('/') ? cleaned : "/#{cleaned}"
  end

  def build_link_map(site)
    site.pages.sort_by(&:path).each_with_object({}) do |page, map|
      next unless markdown_source_page?(site, page)

      export_url = llm_markdown_url(page.url)
      page.data['llm_markdown_url'] = export_url

      canonical = normalize_site_path(page.url)
      map[canonical] = export_url

      if canonical == '/'
        map['/index.html'] = export_url
      elsif canonical.end_with?('/')
        map[canonical.chomp('/')] = export_url
        map[normalize_site_path(File.join(canonical, 'index.html'))] = export_url
      elsif canonical.end_with?('.html')
        map[canonical.sub(/\.html\z/, '')] = export_url
      end
    end
  end

  def markdown_pages(site)
    site.pages.select { |page| markdown_source_page?(site, page) }.sort_by(&:path)
  end

  def strip_front_matter(content)
    content.sub(/\A---\s*\r?\n.*?\r?\n(?:---|\.\.\.)\s*\r?\n/m, '')
  end

  def sanitize_markdown(content)
    sanitized = content.lines.filter_map do |line|
      match = line.match(STANDALONE_ATTR_LIST_PATTERN)
      next line unless match

      anchor_id = match[1].split.find { |token| token.start_with?('#') }
      next nil unless anchor_id

      %(<a id="#{anchor_id.delete_prefix('#')}"></a>\n)
    end.join

    sanitized.gsub(INLINE_ATTR_LIST_PATTERN, '')
  end

  def rewrite_markdown_links(content, page_url, link_map)
    content.gsub(MARKDOWN_LINK_PATTERN) do |match|
      destination = Regexp.last_match[:destination]
      rewritten = rewrite_destination(destination, page_url, link_map)
      next match if rewritten == destination

      match.sub(destination, rewritten)
    end
  end

  def rewrite_destination(destination, page_url, link_map)
    return destination if external_destination?(destination)

    path, query, fragment = destination.match(/\A([^?#]*)(\?[^#]*)?(#.*)?\z/).captures
    return destination if path.nil? || path.empty?

    resolved_path = resolve_site_path(page_url, path)
    rewritten = link_map[resolved_path]
    return destination unless rewritten

    "#{rewritten}#{query}#{fragment}"
  end

  def external_destination?(destination)
    destination.start_with?('#') || destination.match?(%r{\A(?:[a-z][a-z0-9+.-]*:|//)}i)
  end

  def resolve_site_path(page_url, target_path)
    return normalize_site_path(target_path) if target_path.start_with?('/')

    normalize_site_path(File.join(base_dir(page_url), target_path))
  end

  def base_dir(page_url)
    return '/' if page_url == '/'

    page_url.end_with?('/') ? page_url : File.dirname(page_url)
  end

  def export_markdown(site)
    link_map = build_link_map(site)

    markdown_pages(site).each do |page|
      content = File.read(source_path(site, page), encoding: 'UTF-8')
      content = strip_front_matter(content)
      content = sanitize_markdown(content)
      content = rewrite_markdown_links(content, page.url, link_map)
      content = "#{content.rstrip}\n"

      write_output(site, page.data.fetch('llm_markdown_url'), content)
    end
  end

  def export_llms_txt(site)
    config = llms_txt_config(site)
    return unless config['enabled']

    sections = llms_sections(site, config)
    return if sections.empty?

    lines = []

    title = config.fetch('title', site.config['title']).to_s.strip
    unless title.empty?
      lines << "# #{title}"
      lines << ''
    end

    summary = config.fetch('summary', site.config['description']).to_s.strip
    unless summary.empty?
      lines << "> #{summary}"
      lines << ''
    end

    intro = config.fetch('intro', DEFAULT_LLM_INTRO).to_s.strip
    unless intro.empty?
      lines << intro
      lines << ''
    end

    sections.each do |section|
      section_lines = render_section(site, section)
      next if section_lines.empty?

      lines.concat(section_lines)
    end

    write_output(site, config.fetch('path', '/llms.txt'), "#{lines.join("\n").rstrip}\n")
  end

  def llms_sections(site, config)
    sections = Array(config['sections'])
    return sections unless sections.empty?

    auto_sections(site, config)
  end

  def auto_sections(site, config)
    root_title = root_section_title(config)
    grouped_pages = markdown_pages(site)
                    .reject { |page| page.data['llm_exclude'] }
                    .group_by do |page|
                      auto_section_name(
                        page, root_title
                      )
    end

    grouped_pages.map do |title, pages|
      {
        'title' => title,
        '_pages' => dedupe_pages(sort_auto_pages(pages), 'title')
      }
    end.sort_by do |section|
      [section['title'] == root_title ? 0 : 1, section['title'].to_s.downcase]
    end
  end

  def auto_section_name(page, root_title)
    parent = page.data['parent'].to_s.strip
    return root_title if parent.empty?

    parent
  end

  def sort_auto_pages(pages)
    pages.sort_by do |page|
      [page.url == '/' ? 0 : 1, nav_order_value(page), page_title(page).downcase, page.url]
    end
  end

  def nav_order_value(page)
    Integer(page.data['nav_order'])
  rescue StandardError
    Float::INFINITY
  end

  def render_section(site, section)
    entries = render_section_entries(site, section)
    return [] if entries.empty?

    lines = []
    heading = section['title'].to_s.strip
    unless heading.empty?
      lines << "## #{heading}"
      lines << ''
    end

    lines.concat(entries)
    lines << ''
  end

  def render_section_entries(site, section)
    auto_pages = Array(section['_pages'])
    unless auto_pages.empty?
      return auto_pages.filter_map do |page|
        render_page_entry(site, page, section['pages'] || {})
      end
    end

    links = Array(section['links']).filter_map do |link|
      render_link_entry(site, link)
    end
    return links unless links.empty?

    page_config = section['pages']
    return [] unless page_config.is_a?(Hash)

    section_pages(site, page_config).filter_map do |page|
      render_page_entry(site, page, page_config)
    end
  end

  def render_link_entry(site, link)
    title = link.fetch('title', '').to_s.strip
    url = resolve_link_url(site, link)
    return nil if title.empty? || url.empty?

    description = link['description'].to_s.strip
    return "- [#{title}](#{url})" if description.empty?

    "- [#{title}](#{url}): #{description}"
  end

  def resolve_link_url(site, link)
    external_url = link['url'].to_s.strip
    return external_url unless external_url.empty?

    path = link['path'].to_s.strip
    return '' if path.empty?

    absolute_url(site, path)
  end

  def section_pages(site, page_config)
    pages = markdown_pages(site)
    pages = pages.select { |page| page_matches?(page, page_config['where'] || {}) }
    pages = dedupe_pages(pages, page_config['dedupe_by']) if page_config['dedupe_by']

    sort_key = page_config['sort_by']
    return pages unless sort_key

    pages.sort_by { |page| page_value(page, sort_key).to_s.downcase }
  end

  def page_matches?(page, filters)
    filters.all? do |key, expected|
      values_match?(page_value(page, key), expected)
    end
  end

  def values_match?(value, expected)
    return Array(expected).any? { |item| values_match?(value, item) } if expected.is_a?(Array)

    value.to_s == expected.to_s
  end

  def dedupe_pages(pages, dedupe_by)
    keys = Array(dedupe_by).map(&:to_s)
    seen = {}

    pages.each_with_object([]) do |page, result|
      identity = keys.map { |key| page_value(page, key).to_s }
      next if seen[identity]

      seen[identity] = true
      result << page
    end
  end

  def render_page_entry(site, page, page_config)
    label = render_page_template(page_config.fetch('label', '%{title}'), page).strip
    url = absolute_url(site, page.data.fetch('llm_markdown_url'))
    return nil if label.empty? || url.empty?

    description = render_page_template(page_config['description'].to_s, page).strip
    return "- [#{label}](#{url})" if description.empty?

    "- [#{label}](#{url}): #{description}"
  end

  def render_page_template(template, page)
    template.to_s.gsub(/%\{([^}]+)\}/) do
      page_value(page, Regexp.last_match(1)).to_s
    end
  end

  def page_value(page, key)
    case key.to_s
    when 'title'
      page_title(page)
    when 'description'
      page_description(page)
    when 'url'
      page.url
    when 'path'
      page.path
    when 'llm_markdown_url'
      page.data['llm_markdown_url']
    else
      page.data[key.to_s]
    end
  end

  def page_title(page)
    title = page.data['title'].to_s.strip
    return title unless title.empty?

    path = page.path.to_s
    basename = File.basename(path, File.extname(path))
    basename = File.basename(File.dirname(path)) if basename == 'index'
    basename.split(/[-_]/).map(&:capitalize).join(' ')
  end

  def page_description(page)
    llm_description = page.data['llm_description'].to_s.strip
    return llm_description unless llm_description.empty?

    page.data['description'].to_s.strip
  end

  def absolute_url(site, path)
    root = site.config.fetch('url', '').to_s.sub(%r{/*\z}, '')
    baseurl = site.config.fetch('baseurl', '').to_s

    "#{root}#{normalize_site_path("#{baseurl}/#{path}")}"
  end

  def write_output(site, url, content)
    output_path = File.join(site.dest, url.delete_prefix('/'))
    FileUtils.mkdir_p(File.dirname(output_path))
    File.write(output_path, content)
  end

  class Generator < Jekyll::Generator
    safe true
    priority :low

    def generate(site)
      MarkdownForLlms.build_link_map(site)
    end
  end
end

Jekyll::Hooks.register :site, :post_write do |site|
  MarkdownForLlms.export_markdown(site)
  MarkdownForLlms.export_llms_txt(site)
end
