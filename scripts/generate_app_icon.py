#!/usr/bin/env python3
"""Generate Clientele app icon: black background with white ring."""

from pathlib import Path

from PIL import Image, ImageDraw

SIZE = 1024
SCRIPT_DIR = Path(__file__).resolve().parent
OUTPUT = SCRIPT_DIR.parent / "frontend/Clientele Manager/Clientele Manager/Assets.xcassets/AppIcon.appiconset/AppIcon.png"

# Create black canvas
img = Image.new("RGB", (SIZE, SIZE), color="#000000")
draw = ImageDraw.Draw(img)

# Draw white ring: centered circle, inner radius ~350px, outer ~430px (80px stroke)
center = SIZE // 2
inner_radius = 350
outer_radius = 430

# Draw ring by drawing outer circle filled white, then inner circle filled black
draw.ellipse(
    (center - outer_radius, center - outer_radius, center + outer_radius, center + outer_radius),
    fill="#FFFFFF",
)
draw.ellipse(
    (center - inner_radius, center - inner_radius, center + inner_radius, center + inner_radius),
    fill="#000000",
)

OUTPUT.parent.mkdir(parents=True, exist_ok=True)
img.save(OUTPUT, "PNG")
print(f"Generated {OUTPUT} ({SIZE}x{SIZE})")
