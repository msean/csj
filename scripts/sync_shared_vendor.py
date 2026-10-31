#!/usr/bin/env python3
"""Sync the local shared module into the repository's existing Go vendor trees.
Run after editing shared/; --check verifies freshness without modifying files.
Only the shared module and its existing ULID dependency are managed here.
"""
from pathlib import Path
import argparse
import shutil

ROOT = Path(__file__).resolve().parents[1]
MODULE = 'github.com/msean/csj/shared'
PACKAGES = ['constants', 'internal/modelutil', 'model', 'sharding']
parser = argparse.ArgumentParser()
parser.add_argument('--check', action='store_true')
args = parser.parse_args()
for module in ['app', 'backend']:
    vendor = ROOT / module / 'vendor'
    dest = vendor / MODULE
    source_files = {p.relative_to(ROOT / 'shared'): p.read_bytes() for package in PACKAGES for p in (ROOT / 'shared' / package).glob('*.go') if not p.name.endswith('_test.go')}
    actual = {p.relative_to(dest): p.read_bytes() for p in dest.rglob('*.go')} if dest.exists() else {}
    if args.check:
        if source_files != actual:
            raise SystemExit(module + ': shared vendor snapshot is stale; run scripts/sync_shared_vendor.py')
    else:
        if dest.exists(): shutil.rmtree(dest)
        for path, data in source_files.items():
            target = dest / path
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(data)
    manifest = vendor / 'modules.txt'
    text = manifest.read_text()
    header = '# ' + MODULE + ' v0.0.0 => ../shared'
    if header not in text:
        if args.check: raise SystemExit(module + ': missing shared vendor manifest entry')
        text += '\n' + header + '\n## explicit; go 1.24.6\n' + ''.join(MODULE + '/' + package + '\n' for package in PACKAGES) + '# ' + MODULE + ' => ../shared\n'
    if module == 'backend':
        ulid = vendor / 'github.com/oklog/ulid/v2'
        upstream = ROOT / 'app/vendor/github.com/oklog/ulid/v2'
        if args.check:
            if not ulid.exists(): raise SystemExit('backend: missing vendored ULID')
        else:
            if ulid.exists(): shutil.rmtree(ulid)
            shutil.copytree(upstream, ulid)
        if '# github.com/oklog/ulid/v2 v2.1.1' not in text:
            if args.check: raise SystemExit('backend: missing ULID manifest entry')
            text += '\n# github.com/oklog/ulid/v2 v2.1.1\n## explicit; go 1.15\ngithub.com/oklog/ulid/v2\n'
    if not args.check: manifest.write_text(text)
print('Shared vendor snapshots are current.' if args.check else 'Shared vendor snapshots updated.')
