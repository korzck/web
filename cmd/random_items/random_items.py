import csv
import random

first = [
    "Точечное", 
    "Токарное",
    "Термическое",
    "Поверхностное",
    "Фрезерное",
    "Наружное",
    "Шлифовальное",
    "Автоматизированное",
    "Многоточечное",
    "Лазерное",
    "Полировочное",
    "Плазменное"
]

second = [
    "сверление",
    "обрабатывание",
    "отделывание",
    "тиснение",
    "точение",
    "прессование",
    "сгибание",
    "строгание",
    "протяжение",
    "прокатывание",
    "профилирование",
    "литьё",

]

third = [
    "болванок",
    "деталей",
    "листов",
    "заготовок",
    "проволки из",
    "сплавов"
]

forth = [
    "хрома",
    "стали",
    "меди",
    "латуни",
    "алюминия",
    "олова",
    "никеля",
    "железа",
    "вольфрама",
    "бронзы",
    "свинца",
    "титана",
]

imgs = [
    "http://127.0.0.1:9000/cnc/3cd1e53c-fda6-4d4f-bc16-6cf888f0c57c.jpeg",
    "http://127.0.0.1:9000/cnc/7717cb2c-e4a7-42a0-bbd0-f119589d062a.jpeg",
    "http://127.0.0.1:9000/cnc/c742f6be-7661-4f4f-93f2-a88701e695ba.jpg",
    "http://127.0.0.1:9000/cnc/2a30cd91-6227-41b8-b108-10610de83f68.jpg",
    "http://127.0.0.1:9000/cnc/44975b37-d33a-4574-bec2-c61c5d54985a.jpeg",
    "http://127.0.0.1:9000/cnc/0d5acddd-8063-4f71-b6d8-d4d72d167671.jpeg"
]

c = 0
l = []

for f in first:
    for s in second:
        for t in third:
            for fo in forth:
                for fi in forth:
                    if fi == fo:
                        continue
                    c += 1
                    typeT = "metal"
                    if t == "сплавов":
                        typeT = "alloy"
                    if t == "проволки из":
                        typeT = "wire"
                    l.append({
                        "title": '{f} {s} {t} {fo} и {fi}'.format(f=f, s=s, t=t, fo=fo, fi=fi),
                        "type": typeT,
                        "subtitle": '{f} {s} {t} {fo} и {fi}'.format(f=f, s=s, t=t, fo=fo, fi=fi),
                        "price": random.randint(1, 1000),
                        "imgurl": random.choice(imgs),
                        "info": '{f} {s} {t} {fo} и {fi}'.format(f=f, s=s, t=t, fo=fo, fi=fi)
                    })
                if t == "сплавов":
                    continue
                typeT = "clean"
                if t == "проволки из":
                    typeT = "wire"
                l.append({
                        "title": '{f} {s} {t} {fo}'.format(f=f, s=s, t=t, fo=fo),
                        "type": typeT,
                        "subtitle": '{f} {s} {t} {fo}'.format(f=f, s=s, t=t, fo=fo),
                        "price": random.randint(1, 1000),
                        "imgurl": random.choice(imgs),
                        "info": '{f} {s} {t} {fo}'.format(f=f, s=s, t=t, fo=fo)
                    })

random.shuffle(l)
print(len(l))
keys = l[0].keys()

with open('items.csv', 'w', newline='') as output_file:
    dict_writer = csv.DictWriter(output_file, fieldnames=keys)
    dict_writer.writeheader()
    dict_writer.writerows(l)



