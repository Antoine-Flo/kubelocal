# Create a GitHub release
release version:
	git tag {{version}}
	git push origin {{version}}
	git push
