"""Check the generated recovery pages, not just the redirect configuration."""
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import unquote, urljoin, urlparse
import json
import sys

from mkdocs.config import load_config


class Head(HTMLParser):
    def __init__(self):
        super().__init__()
        self.canonical = None
        self.refresh = None
        self.description = None
        self.links = []

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == "link" and attrs.get("rel") == "canonical":
            self.canonical = attrs.get("href")
        if tag == "meta" and attrs.get("http-equiv", "").lower() == "refresh":
            self.refresh = attrs.get("content")
        if tag == "meta" and attrs.get("name") == "description":
            self.description = attrs.get("content")
        if tag == "a" and attrs.get("href"):
            self.links.append(attrs["href"])


def output_path(markdown):
    path = Path(markdown)
    return path.parent / "index.html" if path.name == "index.md" else path.with_suffix("") / "index.html"


def read_head(path):
    head = Head()
    head.feed(path.read_text(encoding="utf-8"))
    return head


def main():
    config = load_config("mkdocs.yml")
    site = Path(config.site_dir)
    mappings = config.plugins["redirects"].config["redirect_maps"]
    failures = []
    checked = 0
    for old, new in mappings.items():
        source = site / output_path(old)
        target = site / output_path(new)
        if not source.is_file() or not target.is_file():
            failures.append(f"Missing recovery output: {old} -> {new}")
            continue
        head = read_head(source)
        if not head.canonical or not head.refresh:
            failures.append(f"Missing canonical/refresh: {old}")
            continue
        resolved = (source.parent / unquote(urlparse(head.canonical).path) / "index.html").resolve()
        if resolved != target.resolve():
            failures.append(f"Wrong redirect destination: {old} -> {head.canonical}")
        target_head = read_head(target)
        if target_head.refresh or not target_head.canonical:
            failures.append(f"Redirect chain or missing target canonical: {new}")
        checked += 1

    descriptions = []
    for restored in ("Security_Secrets.md", "Troubleshooting_FAQ.md"):
        output = site / output_path(restored)
        head = read_head(output)
        if head.refresh or not head.description or not head.canonical:
            failures.append(f"Restored guide metadata missing: {restored}")
        descriptions.append(head.description)
        base = "https://mt4go.mrpc.pro/" + output_path(restored).parent.as_posix() + "/"
        for link in head.links:
            url = urlparse(urljoin(base, link))
            if url.netloc != "mt4go.mrpc.pro":
                continue
            path = site / unquote(url.path).lstrip("/")
            if path.is_dir():
                path = path / "index.html"
            if not path.is_file():
                failures.append(f"Broken rendered guide link: {restored}: {link}")
    if len(set(descriptions)) != 2:
        failures.append("Restored guides repeat the same description")
    if checked != 61:
        failures.append(f"Expected 61 moved-page redirects, found {checked}")
    report = {"redirects": checked, "restored_guides": 2, "failures": failures}
    print(json.dumps(report, indent=2))
    return bool(failures)


if __name__ == "__main__":
    sys.exit(main())
