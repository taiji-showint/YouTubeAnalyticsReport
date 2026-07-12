#!/usr/bin/env python3
"""
Generate 2-column PDF report using reportlab + matplotlib + PIL
Includes graphs, thumbnails, and detailed analytics with Japanese text support
Usage: python3 generate_pdf.py <input_json> <output_pdf>
"""

import json
import sys
import os
import base64
from datetime import datetime
from io import BytesIO
import matplotlib.pyplot as plt
from PIL import Image, ImageDraw, ImageFont
from reportlab.lib.pagesizes import A4, landscape
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import cm
from reportlab.platypus import SimpleDocTemplate, Table, TableStyle, Paragraph, Spacer, PageBreak, Image as RLImage
from reportlab.lib.enums import TA_CENTER, TA_LEFT
from reportlab.lib import colors

def load_json_data(json_file):
    """Load video statistics from JSON file"""
    with open(json_file, 'r', encoding='utf-8') as f:
        return json.load(f)

def generate_graph(daily_stats, video_title):
    """Generate view count graph using matplotlib"""
    if not daily_stats or len(daily_stats) == 0:
        return None

    try:
        # Extract dates and views
        dates = [stat['date'] for stat in daily_stats]
        views = [stat['views'] for stat in daily_stats]

        # Create figure
        fig, ax = plt.subplots(figsize=(4, 2.5), dpi=100)
        ax.plot(range(len(dates)), views, marker='o', linewidth=1.5, markersize=3, color='#1f77b4')
        ax.fill_between(range(len(dates)), views, alpha=0.3, color='#1f77b4')

        # Format
        ax.set_ylabel('Views', fontsize=8)
        ax.tick_params(axis='both', labelsize=7)
        ax.grid(True, alpha=0.3)
        plt.tight_layout()

        # Convert to image file
        temp_path = f"/tmp/graph_{video_title[:10]}.png"
        plt.savefig(temp_path, format='png', dpi=100, bbox_inches='tight')
        plt.close()

        return temp_path
    except Exception as e:
        print(f"Error generating graph: {e}")
        return None

def create_japanese_text_image(text, width=400, height=30, fontsize=10):
    """Create image with Japanese text using PIL"""
    try:
        # Create image
        img = Image.new('RGB', (width, height), color='white')
        draw = ImageDraw.Draw(img)

        # Try to use system font
        font_paths = [
            "/System/Library/Fonts/Hiragino Sans W3.ttc",
            "/Library/Fonts/Hiragino Sans W3.ttc",
            "/System/Library/Fonts/STHeiti Medium.ttc",
        ]

        font = None
        for font_path in font_paths:
            if os.path.exists(font_path):
                try:
                    font = ImageFont.truetype(font_path, fontsize)
                    break
                except:
                    continue

        if font is None:
            font = ImageFont.load_default()

        # Draw text
        draw.text((5, 5), text, fill='#333333', font=font)

        return img
    except Exception as e:
        print(f"Error creating Japanese text image: {e}")
        return None

def create_styles():
    """Create custom styles for the PDF"""
    styles = getSampleStyleSheet()

    title_style = ParagraphStyle(
        'CustomTitle',
        fontSize=16,
        textColor='#333333',
        spaceAfter=10,
        alignment=TA_LEFT,
        fontName='Helvetica-Bold'
    )

    small_style = ParagraphStyle(
        'SmallText',
        fontSize=8,
        leading=9,
        textColor='#333333',
        spaceAfter=2,
        fontName='Helvetica'
    )

    return {
        'title': title_style,
        'small': small_style,
    }

def format_number(num):
    """Format number with thousands separator"""
    if isinstance(num, float):
        if num == int(num):
            return f"{int(num):,}"
        return f"{num:,.1f}"
    return f"{int(num):,}"

def create_video_info_cell(video, styles, image_dir='reports/images'):
    """Create cell with video information"""
    elements = []

    # Thumbnail
    thumb_path = f"{image_dir}/thumbnail_{video.get('Video_id', 'N/A')}_trim.jpg"
    if os.path.exists(thumb_path):
        try:
            img = RLImage(thumb_path, width=3.5*cm, height=2*cm)
            elements.append(img)
        except:
            elements.append(Paragraph("No image", styles['small']))
    else:
        elements.append(Paragraph("No image", styles['small']))

    elements.append(Spacer(1, 0.15*cm))

    # Video title
    title = video.get('Video_title', 'N/A')[:40]
    elements.append(Paragraph(f"<b>{title}</b>", styles['small']))

    # Video ID, date
    video_id = video.get('Video_id', 'N/A')
    pub_date = video.get('Updated_date', 'N/A')
    duration = int(video.get('Duration', 0))
    elements.append(Paragraph(f"ID: {video_id}", styles['small']))
    elements.append(Paragraph(f"公開: {pub_date} ({duration}日)", styles['small']))

    elements.append(Spacer(1, 0.1*cm))

    # Key metrics
    views = format_number(video.get('View_counts', 0))
    likes = format_number(video.get('Like_counts', 0))
    impressions = format_number(video.get('Impressions', 0))
    ctr = video.get('CTR', 0)

    elements.append(Paragraph(f"<b>再生:</b> {views} | <b>いいね:</b> {likes}", styles['small']))
    elements.append(Paragraph(f"<b>インプレッション:</b> {impressions}", styles['small']))
    elements.append(Paragraph(f"<b>クリック率:</b> {ctr:.2f}%", styles['small']))

    elements.append(Spacer(1, 0.1*cm))

    # Graph
    daily_stats = video.get('daily_stats', [])
    if daily_stats:
        graph_path = generate_graph(daily_stats, video.get('Video_title', 'video'))
        if graph_path and os.path.exists(graph_path):
            try:
                graph_img = RLImage(graph_path, width=4*cm, height=2.5*cm)
                elements.append(graph_img)
            except:
                pass

    elements.append(Spacer(1, 0.1*cm))

    # Age/Gender
    age = video.get('Age_percentage', {})
    gender = video.get('Gender_percentage', {})
    elements.append(Paragraph(
        f"年齢: 25-34 ({age.get('AGE25_34', 0):.0f}%) 35-44 ({age.get('AGE35_44', 0):.0f}%) | "
        f"性別: 男性 ({gender.get('MALE', 0):.0f}%)",
        styles['small']
    ))

    # Traffic
    traffic = video.get('Traffic_source', {})
    elements.append(Paragraph(
        f"流入: 登録者 ({traffic.get('SUBSCRIBER', 0):.0f}%) 関連 ({traffic.get('RELATED_VIDEO', 0):.0f}%) "
        f"検索 ({traffic.get('YT_SEARCH', 0):.0f}%)",
        styles['small']
    ))

    # External sites
    ext_sites = video.get('External_sites', [])
    if ext_sites:
        ext_text = ", ".join([f"{list(site.keys())[0]}" for site in ext_sites[:2]])
        elements.append(Paragraph(f"外部: {ext_text}", styles['small']))

    return elements

def generate_pdf(json_file, pdf_file):
    """Generate PDF from JSON data"""

    print("Loading data...")
    videos = load_json_data(json_file)
    styles = create_styles()

    # Create PDF
    doc = SimpleDocTemplate(
        pdf_file,
        pagesize=landscape(A4),
        rightMargin=0.8*cm,
        leftMargin=0.8*cm,
        topMargin=1*cm,
        bottomMargin=0.8*cm
    )

    story = []

    # Title page
    today = datetime.now().strftime('%Y年%m月%d日')
    title = Paragraph(f"show int レポート {today}", styles['title'])
    story.append(title)
    story.append(Spacer(1, 0.3*cm))

    # Video list
    video_names = []
    for i, video in enumerate(videos, 1):
        video_names.append(f"{i}. {video.get('Video_title', 'N/A')[:35]}")

    video_list_text = " | ".join([f"{i+1}" for i in range(len(videos))])
    video_list = Paragraph(f"動画数: {len(videos)}件", styles['small'])
    story.append(video_list)
    story.append(PageBreak())

    # 2-column video details
    col_width = 9*cm

    print("Generating video cards...")
    for i in range(0, len(videos), 2):
        left_elements = create_video_info_cell(videos[i], styles)

        right_elements = []
        if i + 1 < len(videos):
            right_elements = create_video_info_cell(videos[i + 1], styles)

        row_data = [left_elements, right_elements if right_elements else [""]]

        table = Table(row_data, colWidths=[col_width, col_width])
        table.setStyle(TableStyle([
            ('VALIGN', (0, 0), (-1, -1), 'TOP'),
            ('LEFTPADDING', (0, 0), (-1, -1), 0.2*cm),
            ('RIGHTPADDING', (0, 0), (-1, -1), 0.2*cm),
            ('TOPPADDING', (0, 0), (-1, -1), 0.2*cm),
            ('BOTTOMPADDING', (0, 0), (-1, -1), 0.3*cm),
            ('BORDER', (0, 0), (-1, -1), 0.5, colors.HexColor('#dddddd')),
            ('BACKGROUND', (0, 0), (-1, -1), colors.HexColor('#fafafa')),
        ]))

        story.append(table)
        story.append(Spacer(1, 0.3*cm))

    # Build PDF
    print("Building PDF...")
    doc.build(story)
    print(f"✓ PDF generated: {pdf_file}")

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
        print(f"Error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
