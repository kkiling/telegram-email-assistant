import imgkit
import sys
import os
from PIL import Image
from threading import Thread

directory = sys.argv[1:][0]
have_cid = sys.argv[1:][1] == "true"

htmlPath = os.path.join(directory, "index.html")
pngPath = os.path.join(directory, "index.png")

try:
    if have_cid:
        options = {
            'format': 'png',
            'encoding': 'utf-8',
            'enable-local-file-access': None,
        }
        thread = Thread(target=imgkit.from_file, args=(htmlPath, pngPath, options,))
        thread.start()
        thread.join(60)

    # Try without local files
    if not os.path.exists(pngPath):
        options = {
            'format': 'png',
            'encoding': 'utf-8',
        }
        imgkit.from_file(htmlPath, pngPath, options,)
    
except Exception as e: 
    print(e)

# Optimizing image size
if os.path.exists(pngPath):
    picture = Image.open(pngPath)
    picture.save(pngPath)