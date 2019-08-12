Quick notes:

# Library API

[x] glicense.All()
[x] glicense.AllInfo()
[x] glicense.GetByInfo()
[x] glicense.SearchByName()

[ ] glicense.Detect("MIT License Copyright (c) Permission is hereby granted...")
[ ] glicense.DetectFromPath("/path/to/source/of/license/file")
[ ] glicense.DetectFromURL("https://github.com/abc/")

[ ] glicense.Add("MIT License Copyright...", "/path/to/source")
[ ] glicense.AddWithOption("MIT License Copyright...", "/path/to/source", glicense{
	ExcludedPattern: "*.sql",
	IncludedPattern: "*.go",
	MappingComment: glicense.Mapping {
		".go": "/*{content}*/",
		".xml": "<!--{content}-->",
	},
})

# Commandline

[ ] glicense detect "MIT License Copyright (c) Permission is hereby granted..."
[ ] glicense detect -p /path/to/source/
[ ] glicense detect -u https://github.com/abc/
[ ] echo "MIT License Copyright (c) Permission is hereby granted..." | glicense detect

[ ] glicense add "MIT License Copyright (c) Permission is hereby granted..." /path/to/source/
[ ] glicense add -f /path/to/license/file /path/to/source/

# REST API

# WebUI
