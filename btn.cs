    class RoundedButton : System.Windows.Forms.Button
    {
        public const int WM_CREATE = 0x0001;
        public const int WM_NCCREATE = 0x0081;
        public const int WM_PAINT = 0x000F;

        private int nRadius = 20;

        public int Radius
        {
            get { return nRadius; }
            set { nRadius = value; }
        }

        private int nBorderSize = 4;

        public int BorderSize
        {
            get { return nBorderSize; }
            set { nBorderSize = value; }
        }

        public RoundedButton(int nRadius = 20, System.Drawing.Color? fillColor = null, System.Drawing.Color? borderColor = null, int borderSize = 4)
        {
            Radius = nRadius;
            FillColor = fillColor ?? System.Drawing.Color.White;
            BorderColor = borderColor ?? System.Drawing.Color.Red;
            BorderSize = borderSize;
        }

        private System.Drawing.Color borderColor = System.Drawing.Color.Red;

        public System.Drawing.Color BorderColor
        {
            get { return borderColor; }
            set { borderColor = value; }
        }

        private System.Drawing.Color fillColor = System.Drawing.Color.White;

        public System.Drawing.Color FillColor
        {
            get { return fillColor; }
            set { fillColor = value; }
        }

        private System.Drawing.Color oldFillColor;

        protected override void OnMouseEnter(EventArgs e)
        {
            oldFillColor = FillColor;
            FillColor = System.Drawing.SystemColors.ButtonHighlight;
            this.Refresh();
            base.OnMouseEnter(e);
        }

        protected override void OnMouseLeave(EventArgs e)
        {
            base.OnMouseLeave(e);
            FillColor = oldFillColor;
            this.Refresh();
        }

        protected override void OnCreateControl()
        {
            base.OnCreateControl();
            this.BackColor = FillColor;
            int nShift = BorderSize;
        }

        protected override void WndProc(ref Message m)
        {
            base.WndProc(ref m);
            if (m.Msg == WM_CREATE)
            {
                using (Graphics gr = Graphics.FromHwnd(Handle))
                {
                    gr.InterpolationMode = System.Drawing.Drawing2D.InterpolationMode.HighQualityBilinear;
                    gr.CompositingQuality = System.Drawing.Drawing2D.CompositingQuality.HighQuality;
                    gr.SmoothingMode = System.Drawing.Drawing2D.SmoothingMode.AntiAlias;
                    using (System.Drawing.Drawing2D.GraphicsPath gp = CreatePath(new Rectangle(System.Drawing.Point.Empty, base.Size), nRadius, false))
                    {
                        gr.FillPath(SystemBrushes.Window, gp);
                        Region region = new Region(gp);
                        base.Region = region;
                    }
                }
                m.Result = (IntPtr)1;
            }
            else if (m.Msg == WM_PAINT)
            {
                using (Graphics gr = Graphics.FromHwnd(m.HWnd))
                {
                    gr.InterpolationMode = System.Drawing.Drawing2D.InterpolationMode.HighQualityBilinear;
                    gr.CompositingQuality = System.Drawing.Drawing2D.CompositingQuality.HighQuality;
                    gr.SmoothingMode = System.Drawing.Drawing2D.SmoothingMode.AntiAlias;
                    using (System.Drawing.Drawing2D.GraphicsPath gp = CreatePath(new Rectangle(System.Drawing.Point.Empty, base.Size), nRadius, true))
                    {
                        System.Drawing.Pen redPen = new System.Drawing.Pen(BorderColor, BorderSize);
                        gr.FillPath(new SolidBrush(FillColor), gp);
                        gr.DrawPath(redPen, gp);
                    }

                    System.Drawing.Size textSize = TextRenderer.MeasureText(this.Text, this.Font);
                    var nWidth = ((this.Width - textSize.Width) / 2);
                    var nHeight = ((this.Height - textSize.Height) / 2);
                    System.Drawing.Point drawPoint = new System.Drawing.Point(nWidth, nHeight);
                    Rectangle normalRect = new Rectangle(drawPoint, textSize);
                    TextRenderer.DrawText(gr, this.Text, this.Font, normalRect, ForeColor);
                }
                m.Result = (IntPtr)0;
            }
        }
        public static System.Drawing.Drawing2D.GraphicsPath CreatePath(Rectangle rect, int nRadius, bool bOutline)
        {
            int nShift = bOutline ? 1 : 0;
            System.Drawing.Drawing2D.GraphicsPath path = new System.Drawing.Drawing2D.GraphicsPath();
            path.AddArc(rect.X + nShift, rect.Y, nRadius, nRadius, 180f, 90f);
            path.AddArc((rect.Right - nRadius) - nShift, rect.Y, nRadius, nRadius, 270f, 90f);
            path.AddArc((rect.Right - nRadius) - nShift, (rect.Bottom - nRadius) - nShift, nRadius, nRadius, 0f, 90f);
            path.AddArc(rect.X + nShift, (rect.Bottom - nRadius) - nShift, nRadius, nRadius, 90f, 90f);
            path.CloseFigure();
            return path;
        }
    }
