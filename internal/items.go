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
			Title:    "Сварка",
			Subtitle: "Сварка всех типов сплавов и металлов",
			Price:    "1",
			ImgURL:   "/static/img/img1.jpeg",
			Info: `
			Наша компания специализируется на предоставлении услуг сварки всех типов сплавов и металлов. Мы обладаем богатым опытом и экспертизой в области сварки, и наша команда квалифицированных специалистов готова взяться за самые сложные проекты.

Мы оснащены современным оборудованием, включая ЧПУ станки, которые обеспечивают высокую точность и качество сварочных работ. ЧПУ (числовое программное управление) позволяет нам создавать сварные конструкции по заданным точным спецификациям и чертежам. Это обеспечивает нам возможность работать с различными типами материалов, включая стали, нержавеющую сталь, алюминий, титан и другие сплавы.

Мы ориентированы на удовлетворение потребностей наших клиентов и стараемся предоставлять услуги сварки высокого качества. Мы тщательно следим за соблюдением всех необходимых стандартов и процедур безопасности, чтобы гарантировать, что наши свар
			`,
		},
		{
			Type:     "metal",
			Title:    "Выточка",
			Subtitle: "Выточка и резка металлов",
			Price:    "2",
			ImgURL:   "/static/img/img2.jpeg",
			Info:     "Metal milling, also known as machining or metal cutting, is a process used to shape and create complex parts or components from solid metal blocks. It involves removing material from the workpiece using rotating cutting tools, such as end mills or drills, to achieve the desired shape, size, and surface finish. The milling process can be performed manually or using computer numerical control (CNC) machines, allowing for precision and repeatability. Metal milling is commonly used in industries such as manufacturing, aerospace, automotive, and construction to produce components with tight tolerances and intricate geometries. It is a versatile and efficient method for creating a wide range of metal parts for various applications.  ",
		},
		{
			Type:     "metal",
			Title:    "3д формовка металла",
			Subtitle: "3д формовка металла с высокоточными сверлами",
			Price:    "3",
			ImgURL:   "/static/img/img3.jpeg",
			Info: `
			Выточка и резка металлов - это процессы обработки материала, которые играют важную роль в области металлообработки и машиностроения. Они позволяют создавать сложные детали и изделия из различных металлических материалов с высокой точностью и качеством.

Выточка металлов является процессом удаления материала из обрабатываемой детали с использованием режущего инструмента, такого как токарный инструмент. Она часто применяется для создания различных валов, втулок, фланцев, резьбовых соединений и других компонентов, требующих точной формы и поверхностей. Выточка позволяет обрабатывать детали разного диаметра, длины и формы, в зависимости от потребностей проекта.

Резка металлов, в свою очередь, представляет собой процесс разделения больших листов или пластин металла на более мелкие детали или заготовки. Существует несколько способов резки металла, включая лазерную резку, плазменную резку и газовую резку. Каждый метод обладает своими преимуществами и применяется в зависимости от толщины и типа металла, а также требований к точности и скорости резки.
			`,
		},
		{
			Type:     "metal",
			Title:    "Гравировка по металлу",
			Subtitle: "Гравировка по металлу и полировка",
			Price:    "4",
			ImgURL:   "/static/img/img4.jpeg",
			Info: `
			Гравировка по металлу и полировка - это два процесса, которые применяются для добавления уникальных дизайнов, маркировки и повышения эстетического качества металлических поверхностей.

Гравировка по металлу - это процесс создания рельефных или фрезерных узоров, текстов или изображений на металлических поверхностях. Она может быть выполнена либо вручную с использованием ручных инструментов и станков, либо с помощью современных компьютеризированных систем гравировки. Гравировка может применяться для кастомизации ювелирных изделий, наружных знаков, металлических плиток, памятников, инструментов и многих других предметов. Она не только придает уникальный вид, но также может использоваться для идентификации и брендирования.

Полировка металла - это процесс придания металлическим поверхностям гладкого блеска и удаления дефектов, таких как царапины, вмятины или окислы. Она может быть выполнена с использованием ручного инструмента, абразивных материалов или специальных полировальных машин. Полировка может применяться для улучшения внешнего вида и отделки металлических изделий, таких как ювелирные украшения, автомобильные детали, металлическая мебель, музыкальные инструменты и другие изделия. Она помогает создать поверхность, отражающую свет, и придает профессиональный и элегантный вид.
			`,
		},
		{
			Type:     "wood",
			Title:    "Формовка дерева и фанеры",
			Subtitle: "Формовка дерева и фанеры для болванок любых размеров",
			Price:    "5",
			ImgURL:   "/static/img/img5.jpeg",
			Info: `
			Наша компания предоставляет услуги по формовке дерева и фанеры для болванок любых размеров. Мы специализируемся на использовании компьютерно-числового управления (ЧПУ) станков, которые обеспечивают высокую точность и повторяемость процесса формовки.

При формовке дерева и фанеры мы используем передовые технологии и инновационные методы, чтобы обеспечить высокий уровень качества и точности каждой болванки. Наши опытные специалисты проводят тщательный анализ и проектирование процесса формовки, чтобы достичь оптимальных результатов.

Мы гордимся тем, что можем предложить нашим клиентам широкий спектр возможностей в формовке дерева и фанеры. Мы работаем с различными видами дерева, включая традиционные породы, такие как дуб, сосна и береза, а также экзотические и специальные виды дерева.
			`,
		},
		{
			Type:     "wood",
			Title:    "Резка дерева и фанеры",
			Subtitle: "Резка по дереву и фанеры для болванок любых размеров",
			Price:    "6",
			ImgURL:   "/static/img/img6.jpeg",
			Info: `
			Резка по дереву и фанеры для болванок любых размеров – это процесс, который позволяет создавать высококачественные детали и заготовки из дерева и фанеры согласно индивидуальным требованиям заказчика.

Резка по дереву и фанеры выполняется с использованием специального оборудования, такого как станки с числовым программным управлением (ЧПУ), которые позволяют максимально точно и эффективно выполнять режущие операции. Следуя заданной программе, режущий инструмент проходит через материал, обеспечивая точные и четкие линии реза.
			`,
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
			<div class="button">Добавить в заказ</div>
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
