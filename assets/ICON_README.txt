Icon Asset Required
===================

The NSIS installer requires an icon.ico file in this directory.

Options to obtain an icon:

1. Create using online tools:
   - https://convertio.co/png-ico/
   - https://www.icoconverter.com/

2. Use existing icon from another application (with permission)

3. Generate programmatically using ImageMagick:
   convert -background none -size 256x256 -gravity center \
   label:"TF" -fill "#004225" icon.png
   convert icon.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico

4. Design professionally using:
   - Adobe Illustrator
   - Figma
   - Canva

Requirements:
- Format: .ico (Windows Icon)
- Sizes: Include 16x16, 32x32, 48x48, 64x64, 128x128, 256x256
- Colors: Use brand colors (British Racing Green #004225 preferred)
- Transparency: Recommended for modern look

Temporary Workaround:
For testing, you can comment out the MUI_ICON lines in the NSIS script.
The installer will use the default NSIS icon.
