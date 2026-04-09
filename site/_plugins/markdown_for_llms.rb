require 'fileutils'
require 'pathname'

module MarkdownForLlms
  extend self

  MARKDOWN_EXTENSIONS = ['.md', '.markdown'].freeze
  MARKDOWN_LINK_PATTERN = /(?<!!)(\[[^\]]+\]\()(?<destination>[^)\s]+)(?<suffix>\))/
  STANDALONE_ATTR_LIST_PATTERN = /^\{:\s*([^}]+)\}\s*$/
  INLINE_ATTR_LIST_PATTERN = /\{:\s*[^}]+\}/
  GUIDES = [
    ['Home', '/index.md', 'Project overview and docs entry point.'],
    ['Getting started', '/start.md', 'Installation notes and platform requirements.'],
    ['CLI', '/cli/index.md', 'CLI overview with shell completion notes.'],
    ['TUI', '/tui/index.md', 'TUI overview and shared key bindings.']
  ].freeze
  OPTIONAL_LINKS = [
    ['GitHub repository', 'https://github.com/skatkov/devtui', 'Source code, issues, and project metadata.'],
    ['GitHub releases', 'https://github.com/skatkov/devtui/releases', 'Binary downloads and release history.']
  ].freeze

  def markdown_source_page?(site, page)
    return false unless MARKDOWN_EXTENSIONS.include?(File.extname(page.path).downcase)

    File.file?(source_path(site, page))
  end

  def source_path(site, page)
    File.join(site.source, page.path)
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
    seen_index_entries = {}

    site.pages.sort_by(&:path).each_with_object({}) do |page, map|
      next unless markdown_source_page?(site, page)

      export_url = llm_markdown_url(page.url)
      page.data['llm_markdown_url'] = export_url

      title = page.data['title'].to_s.strip
      parent = page.data['parent'].to_s.strip
      entry_key = [parent, title]
      page.data.delete('llm_index_entry')
      if !title.empty? && !parent.empty? && !seen_index_entries[entry_key]
        page.data['llm_index_entry'] = true
        seen_index_entries[entry_key] = true
      end

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
    site.pages.select { |page| markdown_source_page?(site, page) }
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
    lines = [
      '# DevTUI',
      '',
      '> DevTUI is an all-in-one terminal toolkit for developers with both a CLI and an interactive TUI for formatting, transforming, and inspecting common developer data.',
      '',
      "Prefer the Markdown URLs below over the HTML pages. Each link points to an LLM-friendly Markdown export generated from the site's source pages.",
      ''
    ]

    append_named_links(lines, 'Guides', GUIDES) do |title, path, description|
      "- [#{title}](#{absolute_url(site, path)}): #{description}"
    end

    append_named_links(lines, 'CLI Commands', llm_index_pages(site, 'CLI')) do |page|
      title = page.data.fetch('title')
      "- [#{title}](#{absolute_url(site,
                                   page.data.fetch('llm_markdown_url'))}): Command reference for `devtui #{title}`."
    end

    append_named_links(lines, 'TUI Tools', llm_index_pages(site, 'TUI')) do |page|
      title = page.data.fetch('title')
      "- [#{title}](#{absolute_url(site,
                                   page.data.fetch('llm_markdown_url'))}): TUI documentation for the #{title} tool."
    end

    append_named_links(lines, 'Optional', OPTIONAL_LINKS) do |title, url, description|
      "- [#{title}](#{url}): #{description}"
    end

    write_output(site, '/llms.txt', "#{lines.join("\n")}\n")
  end

  def append_named_links(lines, heading, entries)
    lines << "## #{heading}"
    lines << ''

    entries.each do |entry|
      lines << yield(*Array(entry))
    end

    lines << ''
  end

  def llm_index_pages(site, parent)
    markdown_pages(site)
      .select { |page| page.data['parent'].to_s == parent && page.data['llm_index_entry'] }
      .sort_by { |page| page.data['title'].to_s.downcase }
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
