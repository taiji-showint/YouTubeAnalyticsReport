#!/usr/bin/env python3
"""
Generate 2-column PDF report from JSON video statistics
Usage: python3 generate_pdf.py <input_json> <output_pdf>
"""

import json
import sys
import os
from datetime import datetime
from reportlab.lib.pagesizes import A4, landscape
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import cm
from reportlab.platypus import SimpleDocTemplate, Table, TableStyle, Paragraph, Spacer, PageBreak, Image
from reportlab.lib.enums import TA_CENTER, TA_LEFT
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont

# Register Japanese fonts
def register_japanese_fonts():
    """Register Japanese fonts for PDF generation"""
    # Try multiple font paths and formats
    font_candidates = [
        # macOS standard fonts
        ("/System/Library/Fonts/Hiragino Sans W3.ttc", "Japanese"),
        ("/Library/Fonts/Hiragino Sans W3.ttc", "Japanese"),
        ("/System/Library/Fonts/STHeiti Medium.ttc", "Japanese"),
        # Try to find individual TTF files from Noto Sans CJK
        ("/Library/Fonts/NotoSansCJK-Regular.ttc", "Japanese"),
        ("/usr/local/opt/font-noto-sans-cjk/share/fonts/opentype/NotoSansCJK-Regular.ttc", "Japanese"),
    ]

    font_registered = False
    for font_path, font_name in font_candidates:
        if os.path.exists(font_path):
            try:
                # Try to register the font
                pdfmetrics.registerFont(TTFont(font_name, font_path))
                print(f"Registered Japanese font from: {font_path}")
                font_registered = True
                break
            except Exception as e:
                print(f"Trying alternative: {font_path} failed: {e}")
                continue

    if not font_registered:
        print("Warning: Could not register Japanese font.")
        print("Attempting to use system fonts...")
        # If TTC doesn't work, try using reportlab's built-in CJK support
        try:
            from reportlab.pdfbase.cidfonts import UniFont, CIDFont
            print("Note: Using reportlab's CID font system")
            font_registered = True
        except:
            print("Warning: Text may not display correctly.")

    return font_registered

def load_json_data(json_file):
    """Load video statistics from JSON file"""
    with open(json_file, 'r', encoding='utf-8') as f:
        return json.load(f)

def create_styles(font_name='Japanese'):
    """Create custom styles for the PDF"""
    styles = getSampleStyleSheet()

    # Determine font to use (fallback to Helvetica if Japanese font not available)
    try:
        pdfmetrics.getFont(font_name)
        body_font = font_name
        bold_font = font_name
    except:
        print(f"Font '{font_name}' not available, falling back to Helvetica")
        body_font = 'Helvetica'
        bold_font = 'Helvetica-Bold'

    # Title style
    title_style = ParagraphStyle(
        'CustomTitle',
        parent=styles['Heading1'],
        fontSize=18,
        textColor='#333333',
        spaceAfter=6,
        alignment=TA_CENTER,
        fontName=bold_font
    )

    # Video info style
    info_style = ParagraphStyle(
        'VideoInfo',
        parent=styles['Normal'],
        fontSize=9,
        leading=11,
        textColor='#333333',
        spaceAfter=3,
        fontName=body_font
    )

    # Label style (for "再生回数" etc)
    label_style = ParagraphStyle(
        'Label',
        parent=styles['Normal'],
        fontSize=8,
        leading=9,
        textColor='#666666',
        spaceAfter=1,
        fontName=bold_font
    )

    return {
        'title': title_style,
        'info': info_style,
        'label': label_style,
        'normal': styles['Normal'],
        'heading': styles['Heading2'],
        'body_font': body_font
    }

def format_number(num):
    """Format number with thousands separator"""
    if isinstance(num, float):
        if num == int(num):
            return f"{int(num):,}"
        else:
            return f"{num:,.2f}"
    return f"{int(num):,}"

def create_video_info_cell(video, styles):
    """Create a formatted cell with video information"""
    elements = []

    # Video title (clickable link)
    title = Paragraph(
        f"<b>{video.get('Video_title', 'N/A')}</b>",
        styles['label']
    )
    elements.append(title)

    # URL
    url_text = f"<font size=7>{video.get('Video_id', 'N/A')}</font>"
    elements.append(Paragraph(url_text, styles['info']))

    # Publish date
    elements.append(Paragraph(f"<b>公開日:</b> {video.get('Updated_date', 'N/A')}", styles['info']))

    # Views
    views = format_number(video.get('View_counts', 0))
    elements.append(Paragraph(f"<b>再生回数:</b> {views} 回", styles['info']))

    # Likes/Dislikes
    likes = format_number(video.get('Like_counts', 0))
    dislikes = format_number(video.get('Dislike_counts', 0))
    elements.append(Paragraph(f"<b>グッド/バッド:</b> {likes}/{dislikes}", styles['info']))

    # Impressions (KEY DATA)
    impressions = format_number(video.get('Impressions', 0))
    elements.append(Paragraph(f"<b>インプレッション数:</b> {impressions} 回", styles['info']))

    # CTR (KEY DATA)
    ctr = video.get('CTR', 0)
    elements.append(Paragraph(f"<b>クリック率:</b> {ctr:.2f}%", styles['info']))

    return elements

def generate_pdf(json_file, pdf_file):
    """Generate 2-column PDF from JSON data"""

    # Register Japanese fonts
    font_registered = register_japanese_fonts()
    font_name = 'Japanese' if font_registered else 'Helvetica'

    # Load data
    videos = load_json_data(json_file)
    styles = create_styles(font_name)

    # Create PDF
    doc = SimpleDocTemplate(
        pdf_file,
        pagesize=landscape(A4),
        rightMargin=1*cm,
        leftMargin=1*cm,
        topMargin=1.5*cm,
        bottomMargin=1*cm
    )

    story = []

    # Title page
    today = datetime.now().strftime('%Y年%m月%d日')
    title = Paragraph(f"show int レポート {today}", styles['title'])
    story.append(title)

    subtitle = Paragraph("動画統計情報", styles['heading'])
    story.append(subtitle)
    story.append(Spacer(1, 0.5*cm))

    # Video list (before detail)
    video_list_text = []
    for i, video in enumerate(videos, 1):
        video_list_text.append(f"{i}. {video.get('Video_title', 'N/A')}")

    video_list_para = Paragraph(
        "<br/>".join(video_list_text),
        styles['info']
    )
    story.append(video_list_para)
    story.append(PageBreak())

    # 2-column video details
    for i in range(0, len(videos), 2):
        row_data = []

        # Left column
        left_video = videos[i]
        left_elements = create_video_info_cell(left_video, styles)

        # Right column
        right_elements = []
        if i + 1 < len(videos):
            right_video = videos[i + 1]
            right_elements = create_video_info_cell(right_video, styles)

        # Create table row
        row_data = [
            left_elements if left_elements else [""],
            right_elements if right_elements else [""]
        ]

        table = Table(row_data, colWidths=[10*cm, 10*cm])
        table.setStyle(TableStyle([
            ('VALIGN', (0, 0), (-1, -1), 'TOP'),
            ('LEFTPADDING', (0, 0), (-1, -1), 0.3*cm),
            ('RIGHTPADDING', (0, 0), (-1, -1), 0.3*cm),
            ('TOPPADDING', (0, 0), (-1, -1), 0.3*cm),
            ('BOTTOMPADDING', (0, 0), (-1, -1), 0.5*cm),
            ('BORDER', (0, 0), (-1, -1), 0.5, '#cccccc'),
            ('BACKGROUND', (0, 0), (-1, -1), '#f9f9f9'),
        ]))

        story.append(table)
        story.append(Spacer(1, 0.5*cm))

    # Build PDF
    doc.build(story)
    print(f"PDF generated: {pdf_file}")

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: python3 generate_pdf.py <input_json> <output_pdf>")
        sys.exit(1)

    json_file = sys.argv[1]
    pdf_file = sys.argv[2]

    if not os.path.exists(json_file):
        print(f"Error: JSON file not found: {json_file}")
        sys.exit(1)

    try:
        generate_pdf(json_file, pdf_file)
    except Exception as e:
        print(f"Error generating PDF: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
