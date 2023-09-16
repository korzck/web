package internal

import (
	"fmt"
	"os"
	"strings"
)

type Item struct {
	Title    string
	Subtitle string
	Price    string
	ImgURL   string
	URL      string
	Info     string
	Type     string
}

func GetItems() []*Item {
	items := []*Item{
		{
			Type:     "metal",
			Title:    "Metal welding",
			Subtitle: "Welding with every type of alloys and forms",
			Price:    "1",
			ImgURL:   "/static/img/img1.jpeg",
			Info:     "Metal welding is a process of joining two or more metallic pieces together using heat, pressure, or a combination of both. It involves melting the surfaces of the metals to be joined and allowing them to cool, creating a strong and durable bond. Welding is widely used in various industries, including manufacturing, construction, automotive, and aerospace, to create structures, repair metal components, and fabricate intricate designs. Different welding techniques, such as arc welding, MIG welding, TIG welding, and spot welding, are employed depending on the specific application and types of metals involved.",
		},
		{
			Type:     "metal",
			Title:    "Metal milling",
			Subtitle: "High speed milling for soft and hard metal forms",
			Price:    "2",
			ImgURL:   "/static/img/img2.jpeg",
			Info:     "Metal milling, also known as machining or metal cutting, is a process used to shape and create complex parts or components from solid metal blocks. It involves removing material from the workpiece using rotating cutting tools, such as end mills or drills, to achieve the desired shape, size, and surface finish. The milling process can be performed manually or using computer numerical control (CNC) machines, allowing for precision and repeatability. Metal milling is commonly used in industries such as manufacturing, aerospace, automotive, and construction to produce components with tight tolerances and intricate geometries. It is a versatile and efficient method for creating a wide range of metal parts for various applications.  ",
		},
		{
			Type:     "metal",
			Title:    "3D metal molding",
			Subtitle: "3D shaped molding with 0.8mm nozzle",
			Price:    "3",
			ImgURL:   "/static/img/img3.jpeg",
			Info: `3D metal molding, also known as metal additive manufacturing or metal 3D printing, is an innovative manufacturing process that involves building three-dimensional metal objects layer by layer from a digital model. Unlike traditional metal molding techniques, which often involve subtractive processes like cutting or milling, 3D metal molding adds material layer by layer, allowing for complex geometries and internal structures that would be challenging or impossible to achieve using other methods.

			In 3D metal molding, a powdered metal material, such as stainless steel, titanium, or aluminum, is selectively melted or fused together by a high-powered laser or electron beam. The laser or electron beam fuses the metal particles in the desired areas, creating a solid, fully dense metal part. This process is repeated layer by layer until the complete object is formed.
			
			3D metal molding offers numerous advantages, including greater design freedom, reduced material waste, faster production times, and the ability to create highly customized or intricate metal parts. It has found applications in various industries, including aerospace, automotive, medical, and jewelry, and continues to advance with new materials and improved technologies.  `,
		},
		{
			Type:     "metal",
			Title:    "Metal engraving",
			Subtitle: "Engraving with water cleaning supplies and high speed polishing",
			Price:    "4",
			ImgURL:   "/static/img/img4.jpeg",
			Info: `Metal engraving is a process of etching or cutting designs, text, or patterns onto metal surfaces. It is a technique used to create permanent and detailed markings on various metals, such as stainless steel, aluminum, brass, or gold. 

			There are different methods of metal engraving, including hand engraving, rotary engraving, and laser engraving. Hand engraving involves using special tools, such as gravers or burins, to manually carve into the metal surface. Rotary engraving uses a rotating cutting tool to engrave the desired design, while laser engraving utilizes a laser beam to vaporize or remove layers of the metal, creating precise and intricate engravings. 
			
			Metal engraving has numerous applications, including jewelry design, personalization of metal items, signage, industrial marking, and artwork. It allows for the creation of intricate and durable designs, adding a decorative or functional element to metal objects. In addition, advancements in laser technology have made metal engraving faster and more precise, further expanding its range of applications.  `,
		},
		{
			Type:     "wood",
			Title:    "Plywood and wood molding",
			Subtitle: "Wood molding for any size pieces",
			Price:    "5",
			ImgURL:   "/static/img/img5.jpeg",
			Info: `Plywood is a type of engineered wood product made from thin layers or plies of wood veneers that are glued together in alternating directions. This cross-layered construction gives plywood its strength, stability, and resistance to warping. Plywood comes in various thicknesses, sizes, and grades, and it is widely used in construction, furniture making, cabinetry, and other woodworking applications.

			Wood molding, also known as trim or moulding, refers to decorative or functional strips of wood used to add aesthetic and architectural detail to various surfaces. It is commonly used around windows, doors, ceilings, and floors to enhance the appearance of interior spaces. Wood molding can come in various profiles and styles, such as baseboard molding, crown molding, chair rail molding, and casing molding. It is typically installed using nails or adhesive, and it can be stained, painted, or left unfinished to suit the desired aesthetic.
			
			Both plywood and wood molding play vital roles in woodworking and construction. Plywood provides a versatile and strong material for structural applications, while wood molding adds visual appeal and character to finished spaces.  `,
		},
		{
			Type:     "wood",
			Title:    "Wood cutting",
			Subtitle: "Wood cutting for any size pieces",
			Price:    "6",
			ImgURL:   "/static/img/img6.jpeg",
			Info: `Wood cutting refers to the process of shaping or dividing wood into desired shapes, sizes, or pieces. It can be done through various methods, depending on the specific requirements and tools available.

			One of the most common methods of wood cutting is sawing. Sawing can be done by hand using a handsaw or powered using circular saws, jigsaws, or table saws. Sawing allows for straight cuts, angled cuts, or curved cuts, depending on the type of saw and technique employed. It is suitable for both rough cutting and precise woodworking.
			
			Another method of wood cutting is woodturning, which involves rotating a piece of wood on a lathe and using cutting tools to shape it symmetrically. Woodturning is commonly used to create items such as table legs, bowls, or decorative pieces that require a rounded or cylindrical shape.
			
			Additionally, there are specialized wood cutting techniques such as scroll sawing, bandsawing, and CNC (Computer Numerical Control) routing. These methods provide precise and intricate cuts, allowing for detailed designs and patterns in woodwork.
			
			Wood cutting serves various purposes, ranging from construction and furniture making to crafting and woodworking hobbies. It is essential to use the appropriate cutting tool and technique for the desired results while considering safety precautions to ensure accurate and safe wood cutting.  `,
		},
	}
	return items
}

func InitPages(items []*Item) error {
	content := `
{{template "header.html"}}
    <div class="title header">{{ .Item.Title }}</div>
	<div class="info">
		<img class="image" src="{{ .Item.ImgURL }}">
		<div class="text">{{ .Item.Info }}</div>
		<div class="process">
			<div class="button">Add to orders</div>
			<br />
			<p>Order details</p>
			<textarea type="textarea" class="details-input"></textarea>
		</div>
	</div>
{{template "footer.html"}}
	`

	for _, v := range items {
		file, err := os.Create(fmt.Sprintf("templates/%v.html", strings.ReplaceAll(v.Title, " ", "-")))
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		defer file.Close()

		// Write the string to the file
		_, err = file.WriteString(content)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return err
		}
	}
	return nil
}
