# GitHub Wiki Setup Instructions

This directory contains all the markdown files for the phev2mqtt GitHub wiki.

## Wiki Pages Created

The following wiki pages have been created:

1. **Home.md** - Wiki home page with overview and quick links
2. **Installation.md** - Complete installation guide (Docker & Unraid)
3. **Configuration.md** - Comprehensive configuration reference
4. **Quick-Start.md** - 5-minute quick start guide
5. **Home-Assistant-Integration.md** - Home Assistant setup and integration
6. **MikroTik-Integration.md** - MikroTik WiFi bridge configuration
7. **WiFi-Management.md** - WiFi management features and configuration
8. **Development.md** - Development environment and contribution guide
9. **Protocol-Documentation.md** - PHEV protocol technical documentation
10. **Troubleshooting.md** - Common issues and solutions
11. **Security-Best-Practices.md** - Security hardening and best practices

## How to Set Up the GitHub Wiki

GitHub wikis are actually separate Git repositories. Here's how to publish these pages:

### Method 1: Using GitHub Web Interface (Easiest)

1. **Enable the wiki on your repository:**
   - Go to your GitHub repository
   - Click **Settings**
   - Scroll to **Features** section
   - Check **Wikis** to enable

2. **Create initial wiki page:**
   - Click the **Wiki** tab
   - Click **Create the first page**
   - This initializes the wiki

3. **Upload pages via web interface:**
   - For each .md file in this directory:
     - Click **New Page**
     - Copy the filename (without .md) as the page title
     - Paste the content from the file
     - Click **Save Page**

### Method 2: Using Git (Recommended for bulk upload)

1. **Enable the wiki** (see Method 1, step 1-2)

2. **Clone the wiki repository:**
   ```bash
   git clone https://github.com/stefanh12/phev2mqtt.wiki.git
   cd phev2mqtt.wiki
   ```

3. **Copy all wiki markdown files:**
   ```bash
   cp /path/to/phev2mqtt/wiki/*.md .
   ```

4. **Commit and push:**
   ```bash
   git add *.md
   git commit -m "Add comprehensive wiki documentation"
   git push origin master
   ```

5. **Verify:**
   - Go to your repository's Wiki tab
   - All pages should now be visible

### Method 3: Using GitHub CLI

1. **Install GitHub CLI** if not already installed:
   ```bash
   # macOS
   brew install gh
   
   # Linux
   sudo apt install gh
   ```

2. **Authenticate:**
   ```bash
   gh auth login
   ```

3. **Clone and push wiki:**
   ```bash
   # Clone wiki
   gh repo clone stefanh12/phev2mqtt.wiki
   cd phev2mqtt.wiki
   
   # Copy files
   cp /path/to/phev2mqtt/wiki/*.md .
   
   # Commit and push
   git add *.md
   git commit -m "Add comprehensive wiki documentation"
   git push origin master
   ```

## Wiki Page Structure

The wiki is organized as follows:

```
Home (Home.md)
├── Getting Started
│   ├── Quick Start (Quick-Start.md)
│   ├── Installation (Installation.md)
│   └── Configuration (Configuration.md)
├── Integration & Features
│   ├── Home Assistant Integration (Home-Assistant-Integration.md)
│   ├── MikroTik Integration (MikroTik-Integration.md)
│   └── WiFi Management (WiFi-Management.md)
└── Advanced Topics
    ├── Development (Development.md)
    ├── Protocol Documentation (Protocol-Documentation.md)
    ├── Security Best Practices (Security-Best-Practices.md)
    └── Troubleshooting (Troubleshooting.md)
```

## Customizing the Sidebar

To create a custom sidebar for easy navigation:

1. **Create _Sidebar.md:**
   ```bash
   cd phev2mqtt.wiki
   nano _Sidebar.md
   ```

2. **Add navigation links:**
   ```markdown
   ## phev2mqtt Wiki
   
   ### Getting Started
   * [Home](Home)
   * [Quick Start](Quick-Start)
   * [Installation](Installation)
   * [Configuration](Configuration)
   
   ### Integration
   * [Home Assistant](Home-Assistant-Integration)
   * [MikroTik](MikroTik-Integration)
   * [WiFi Management](WiFi-Management)
   
   ### Advanced
   * [Development](Development)
   * [Protocol Docs](Protocol-Documentation)
   * [Troubleshooting](Troubleshooting)
   * [Security](Security-Best-Practices)
   ```

3. **Commit and push:**
   ```bash
   git add _Sidebar.md
   git commit -m "Add wiki sidebar"
   git push origin master
   ```

## Customizing the Footer

To add a custom footer to all wiki pages:

1. **Create _Footer.md:**
   ```bash
   cd phev2mqtt.wiki
   nano _Footer.md
   ```

2. **Add footer content:**
   ```markdown
   ---
   📝 Found an error? [Edit this page](https://github.com/stefanh12/phev2mqtt/wiki/_edit) | 
   💬 [Discussions](https://github.com/stefanh12/phev2mqtt/discussions) | 
   🐛 [Report Issue](https://github.com/stefanh12/phev2mqtt/issues)
   ```

3. **Commit and push:**
   ```bash
   git add _Footer.md
   git commit -m "Add wiki footer"
   git push origin master
   ```

## Maintaining the Wiki

### Updating Pages

To update wiki pages:

1. **Clone the wiki repository:**
   ```bash
   git clone https://github.com/stefanh12/phev2mqtt.wiki.git
   cd phev2mqtt.wiki
   ```

2. **Make changes:**
   ```bash
   nano Page-Name.md
   ```

3. **Commit and push:**
   ```bash
   git add Page-Name.md
   git commit -m "Update Page-Name with new information"
   git push origin master
   ```

### Adding New Pages

1. **Create the markdown file:**
   ```bash
   cd phev2mqtt.wiki
   nano New-Page-Name.md
   ```

2. **Add content and commit:**
   ```bash
   git add New-Page-Name.md
   git commit -m "Add New-Page-Name"
   git push origin master
   ```

3. **Link to new page** from relevant existing pages

## Best Practices

### Writing Style

- ✅ Use clear, concise language
- ✅ Include code examples
- ✅ Add screenshots where helpful
- ✅ Link between related pages
- ✅ Keep information up to date

### Formatting

- ✅ Use proper markdown headers (H1, H2, H3)
- ✅ Format code blocks with language hints
- ✅ Use tables for structured data
- ✅ Add horizontal rules to separate sections
- ✅ Include navigation at top and bottom of long pages

### Linking

**Internal wiki links:**
```markdown
[Configuration](Configuration)
[Home Assistant Integration](Home-Assistant-Integration)
```

**Links to repository files:**
```markdown
[Main README](../README.md)
[Security Audit](../SECURITY_AUDIT.md)
```

**External links:**
```markdown
[MikroTik](https://mikrotik.com)
[Home Assistant](https://www.home-assistant.io)
```

## Publishing Checklist

Before publishing the wiki, verify:

- [ ] All markdown files are properly formatted
- [ ] All internal links work correctly
- [ ] Code examples are tested and accurate
- [ ] Screenshots are clear and relevant
- [ ] Sidebar navigation is complete
- [ ] Footer provides helpful links
- [ ] Home page provides good overview
- [ ] Contact/support information is current

## Collaborative Editing

### Allow Community Contributions

By default, GitHub wikis are editable by repository collaborators. To allow broader contributions:

1. **Settings** → **Features** → **Wikis**
2. Configure wiki edit permissions
3. Consider creating a CONTRIBUTING.md with wiki guidelines

### Review Process

For wikis, changes are immediate. Consider:
- Monitoring wiki changes via email notifications
- Regular reviews of recent changes
- Setting up RSS feed for wiki updates

## Troubleshooting Setup

### Wiki not showing up

**Solution:**
- Ensure wiki is enabled in repository settings
- Create at least one page via web interface first
- Then clone and add more pages

### Links broken

**Solution:**
- Use page names without .md extension
- Replace spaces with hyphens in links
- Example: `Home Assistant Integration` → `[Link](Home-Assistant-Integration)`

### Images not loading

**Solution:**
- Upload images to wiki repository
- Use relative paths: `![alt text](images/screenshot.png)`
- Or host images in main repository: `![alt text](../docs/images/screenshot.png)`

## Next Steps

After setting up the wiki:

1. **Update main README.md** to link to wiki:
   ```markdown
   ## Documentation
   
   Comprehensive documentation is available in the [Wiki](https://github.com/stefanh12/phev2mqtt/wiki).
   ```

2. **Add wiki links to repository description**

3. **Announce wiki availability** in:
   - GitHub Discussions
   - Docker Hub description
   - Release notes

4. **Monitor and update** wiki regularly

## Resources

- [GitHub Wiki Documentation](https://docs.github.com/en/communities/documenting-your-project-with-wikis)
- [Markdown Guide](https://www.markdownguide.org/)
- [GitHub Flavored Markdown](https://github.github.com/gfm/)

---

**Questions or Issues?**

If you encounter problems setting up the wiki, please open an issue in the main repository.
