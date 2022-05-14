import imgkit
import sys
import os

options = {
    'format': 'png',
    'encoding': 'utf-8',
    'enable-local-file-access': None,
}

directory = sys.argv[1:][0]
htmlPath = os.path.join(directory, "index.html")
pngPath = os.path.join(directory, "index.png")
imgkit.from_file(htmlPath, pngPath, options=options)
