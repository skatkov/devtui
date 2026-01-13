# Open Graph Implementation Guide

This document describes how Open Graph (OG) tags and images were implemented for this Jekyll site.

## Implementation History

### Key Commits

1. **Initial jekyll-seo-tag setup** (`941008870d953b0b00f4cab916d5b16b355569d1` - Nov 24, 2025)
   ```bash
   git show 941008870d953b0b00f4cab916d5b16b355569d1
   ```
   - Added `jekyll-seo-tag` to plugins in `_config.yml`
   - Added gem to `Gemfile`

2. **Full SEO tag implementation** (`690774d6b9c2f92250ca0b4ed030da258b2e8a2b` - Sep 27, 2025)
   ```bash
   git show 690774d6b9c2f92250ca0b4ed030da258b2e8a2b
   ```
   - Added `jekyll-seo-tag` gem
   - Removed custom meta tag includes (`_includes/meta_tags/sharing.html` and `base.html`)
   - Added `{% seo %}` tag to `_layouts/default.html`
   - Configured author info in `_config.yml`

3. **OG image generation** (`f7a98529efb896e0fe25c08f846d32373deafa4a` - Sep 25, 2025)
   ```bash
   git show f7a98529efb896e0fe25c08f846d32373deafa4a
   ```
   - Added `jekyll-og-image` gem for automatic OG image generation

## Current Implementation

### Gems Used

The site uses two main gems for Open Graph functionality:

1. **jekyll-seo-tag** - Generates SEO meta tags (Twitter cards, Open Graph, etc.)
2. **jekyll-og-image** - Generates custom OG images automatically

### Gemfile Configuration

```ruby
group :jekyll_plugins do
  gem "jekyll-feed"
  gem 'jekyll-compose'
  gem 'jekyll-seo-tag'
  gem 'jekyll-og-image'
end
```

### _config.yml Configuration

```yaml
# Plugins
plugins:
  - jekyll-feed
  - jekyll-seo-tag
  - jekyll-og-image

# Site settings
author:
  fullname: Stanislav Katkov
  github: skatkov
  twitter: 5katkov

# Keep generated OG images
keep_files: ["assets/images/og"]

# OG Image configuration
og_image:
  output_dir: "assets/images/og"
  domain: "skatkov.com"
  canvas:
    width: 1200
    height: 630
  background:
    color: "#ffffff"
  title:
    font_family: "monospace"
    font_size: 60
    color: "#1f2937"
    position: "center"
    max_width: 1000
    line_height: 1.2
  author:
    font_family: "sans-serif"
    font_size: 28
    color: "#6b7280"
  border:
    width: 8
    color: "#1f2937"
```

### Layout Implementation

The SEO tag is included in the head section of the layout:

```html
<!doctype html>
<html lang="en">
  <head>
    {% seo %}
    <!-- other includes -->
  </head>
```

## How to Copy This Approach

### Step 1: Add Required Gems

Add to your `Gemfile`:
```ruby
group :jekyll_plugins do
  gem 'jekyll-seo-tag'
  gem 'jekyll-og-image'
end
```

Run:
```bash
bundle install
```

### Step 2: Configure Plugins

Add to your `_config.yml`:
```yaml
plugins:
  - jekyll-seo-tag
  - jekyll-og-image
```

### Step 3: Configure Site Author

Add to your `_config.yml`:
```yaml
author:
  fullname: Your Name
  github: yourusername
  twitter: yourusername
```

### Step 4: Add OG Image Configuration

Add to your `_config.yml`:
```yaml
keep_files: ["assets/images/og"]

og_image:
  output_dir: "assets/images/og"
  domain: "yourdomain.com"
  canvas:
    width: 1200
    height: 630
  # Customize other settings as needed
```

### Step 5: Update Layouts

Replace custom meta tags with the SEO tag in your `_layouts/default.html`:
```html
<head>
  {% seo %}
</head>
```

### Step 6: Build Your Site

```bash
bundle exec jekyll build
```

The `jekyll-og-image` plugin will automatically generate OG images for your posts and pages based on your configuration.

## Additional Resources

- [jekyll-seo-tag Documentation](https://github.com/jekyll/jekyll-seo-tag)
- [jekyll-og-image Documentation](https://github.com/nhoizey/jekyll-og-image)
- [Open Graph Protocol](https://ogp.me/)
