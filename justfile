# Create a GitHub release
release version:
	git push
	gh release create {{version}} --generate-notes
