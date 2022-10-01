release:
	goreleaser release --rm-dist

test-release:
	goreleaser release --skip-publish --rm-dist
