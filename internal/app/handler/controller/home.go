package controller

import (
	"donbarrigon/new/internal/app/ui/view"
	"donbarrigon/new/internal/utils/handler"
)

func Home(c *handler.Context) {
	data := view.HomePageData{
		SiteName:   "Mi Tienda Online",
		UserName:   "Juan",
		IsLoggedIn: true,
		TotalItems: 50,
		Products: []view.Product{
			{
				ID:          1,
				Title:       "Laptop Gamer Pro",
				Description: "Potente laptop con RTX 4070 y 32GB RAM",
				Price:       1299.99,
				ImageUrl:    "https://picsum.photos/400/300?random=1",
				InStock:     true,
			},
			{
				ID:          2,
				Title:       "Mouse Gaming RGB",
				Description: "Mouse ergonómico con 16000 DPI",
				Price:       79.99,
				ImageUrl:    "https://picsum.photos/400/300?random=2",
				InStock:     true,
			},
			{
				ID:          3,
				Title:       "Teclado Mecánico",
				Description: "Switches Cherry MX Blue, retroiluminado",
				Price:       149.99,
				ImageUrl:    "https://picsum.photos/400/300?random=3",
				InStock:     false,
			},
			{
				ID:          4,
				Title:       "Monitor 4K 27\"",
				Description: "Panel IPS, 144Hz, HDR400",
				Price:       499.99,
				ImageUrl:    "https://picsum.photos/400/300?random=4",
				InStock:     true,
			},
			{
				ID:          5,
				Title:       "Auriculares Inalámbricos",
				Description: "Cancelación de ruido activa, 30h batería",
				Price:       199.99,
				ImageUrl:    "https://picsum.photos/400/300?random=5",
				InStock:     true,
			},
			{
				ID:          6,
				Title:       "Webcam Full HD",
				Description: "1080p 60fps, con micrófono integrado",
				Price:       89.99,
				ImageUrl:    "https://picsum.photos/400/300?random=6",
				InStock:     true,
			},
		},
	}

	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	view.WriteHomePage(c.Writer, data)
}
