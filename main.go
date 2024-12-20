package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Product представляет продукт
type Product struct {
	ID          string
	ImageURL    string
	Name        string
	Description string
	Price       float64
}

// Пример списка продуктов
var products = []Product{
	{ID: "1727088258", ImageURL: "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fi.ytimg.com%2Fvi%2FE3Huy2cdih0%2Fmaxresdefault.jpg&f=1&nofb=1&ipt=8c93c275775a5bdd9a93c346d13054bddb4f5c24b20f51eedd3198dce5ad440e&ipo=images", Name: "Elden Ring", Description: `Золотой порядок сломлен

Восстань, погасшая душа! Междуземье ждёт своего повелителя. Пусть благодать приведёт тебя к Кольцу Элден.

ELDEN RING — ролевая игра в жанре фэнтези от FromSoftware Inc. и BANDAI NAMCO Entertainment Inc. и на сегодняшний день — самый масштабный проект FromSoftware, действие которого разворачивается в мире, полном тайн и опасностей.

В Междуземье, владении королевы Марики Бессмертной, разбилось великое Кольцо Элден — источник силы Древа Эрд. Вскоре отпрыски Марики, полубоги, завладели осколками Кольца Элден — великими рунами. Однако обретённая сила развратила детей королевы, и они развязали войну — Раскол. И отреклась от них Великая Воля.

Теперь же благодать ведёт Погасших, лишённых золотой милости и изгнанных из Междуземья. Живые мертвецы, давным-давно утратившие благодать, отправляйтесь в Междуземье через туманное море и предстаньте перед Кольцом Элден.

Ибо оно ждёт нового повелителя`, Price: 2519},
	{ID: "1727108578", ImageURL: "https://avatars.mds.yandex.net/i?id=5f9cf592c011c5ffb388935e31cd5ef2_l-7973815-images-thumbs&n=13", Name: "Persona 5 Royal", Description: `Главный герой — старшеклассник, который был вынужден переехать в Токио и перевестись в одну из городских школ. Вскоре после этого он видит странный сон: таинственный голос называет его узником судьбы и сообщает, что в недалеком будущем юношу ожидает катастрофа. Теперь, чтобы пройти некий курс «реабилитации», он должен спасать людей от их собственных порочных желаний, примерив маску Призрачного похитителя.`, Price: 2238},
	{ID: "1727127024", ImageURL: "https://digiseller.mycdn.ink/preview/1115001/p1_4384167_9f03e94f.jpg", Name: "Black Myth: Wukong", Description: `Black Myth: Wukong — ролевой боевик по мотивам китайской мифологии. Его сюжет основывается на «Путешествии на Запад», одном из четырёх классических романов китайской литературы. Став Избранным, вы отправитесь в приключение, полное испытаний и чудес, в котором вам предстоит приподнять завесу тайны над великой легендой.`, Price: 4441},
	{ID: "1727127095", ImageURL: "https://korobok.store/upload/iblock/d67/na9245kxy88ke4ley3ce4em14915p4tp.webp", Name: "Baldur's Gate 3", Description: `Соберите отряд и вернитесь в Забытые Королевства. Вас ждет история о дружбе и предательстве, выживании и самопожертвовании, о сладком зове абсолютной власти.

Ваш мозг стал вместилищем для личинки иллитида, и она пробуждает в вас таинственные, пугающие способности. Сопротивляйтесь паразиту и обратите тьму против себя самой – или же безоглядно отдайтесь злу и станьте его воплощением.

Ролевая игра нового поколения в мире Dungeons & Dragons от создателей Divinity: Original Sin 2.`, Price: 2469},
	{ID: "1727127227", ImageURL: "https://avatars.mds.yandex.net/i?id=e542066fc8cfc06acaab48c77cf1e4c6_l-7013372-images-thumbs&n=13", Name: "NieR:Automata", Description: `NieR: Automata tells the story of androids 2B, 9S and A2 and their battle to reclaim the machine-driven dystopia overrun by powerful machines.

Humanity has been driven from the Earth by mechanical beings from another world. In a final effort to take back the planet, the human resistance sends a force of android soldiers to destroy the invaders. Now, a war between machines and androids rages on... A war that could soon unveil a long-forgotten truth of the world.`, Price: 1007},
	{ID: "1727127293", ImageURL: "https://cdn11.bigcommerce.com/s-6kgfzq4siu/images/stencil/1280x1280/products/2271/11135/7bafc5e1b3a974e6765995d44c5ed564__25125.1675955718.jpg?c=1", Name: "NieR Replicant", Description: `NieR Replicant ver.1.22474487139... – обновленная версия NieR Replicant, выпущенной ранее только в Японии.
Эта игра представляет собой приквел шедевральной NieR:Automata. Обновленная версия может похвастаться мастерски отреставрированной графикой, захватывающей сюжетной линией – и не только!

Главный герой – добросердечный юноша из глухой деревеньки. Его сестру Йону поразил смертельный недуг, именуемый «черными буквами», и ради нее он вместе с Белым Гримуаром – причудливой говорящей книгой – отправляется на поиски «запечатанных виршей».`, Price: 1646},
	{ID: "1727127365", ImageURL: "https://wholesgame.com/wp-content/uploads/Helldivers-2-Thumb-4-x-5.jpg", Name: "HELLDIVERS™ 2", Description: `Последняя линия нападения галактики.

Станьте Адским Десантником и сражайтесь за свободу по всему враждебному космосу в динамичном, безумном и беспощадном шутере с видом от третьего лица.
ВАЖНОЕ СООБЩЕНИЕ — ВООРУЖЕННЫЕ СИЛЫ СУПЕР-ЗЕМЛИ
Свобода. Мир. Демократия.
Ваши права граждан Супер-Земли — основы нашей цивилизации.
Самого нашего существования.
Но война продолжается. И вновь всему, что нас окружает, грозит опасность.
Вступите в величайшую армию в истории и сделайте галактику безопасной и свободной.`, Price: 3530},
	{ID: "1727127432", ImageURL: "https://app-time.ru/uploads/games/keyart/2023/12/27122023150938.webp", Name: "Persona 3 Reload", Description: `Новая школа, новые друзья и неожиданное попадание в новую реальность, которая доступна в течение одного часа «между сегодня и завтра». Обретите невероятную силу и раскройте тайны Темного часа, сражайтесь за своих друзей и навсегда оставьте след в их памяти.

Persona 3 Reload — захватывающее современное переосмысление культовой ролевой игры.`, Price: 2569},
	{ID: "1727127502", ImageURL: "https://images.wallpapersden.com/image/download/cult-of-the-lamb-hd-gaming_bWhpa2iUmZqaraWkpJRmamhrrWdlaW0.jpg", Name: "Cult of the Lamb", Description: `В Cult of the Lamb вы окажетесь в роли одержимого ягнёнка, спасённого от гибели жутким незнакомцем. Чтобы вернуть долг своему спасителю, вам придётся найти ему верных последователей. Взращивайте собственный культ в землях лжепророков, совершайте походы по таинственным уголкам леса, объединяйте вокруг себя верных последователей и несите своё слово в массы, чтобы сделать свой культ единственным.
СОБЕРИТЕ СВОЮ ПАСТВУ
Собирайте ресурсы и используйте их для возведения построек, проводите темные ритуалы, чтобы задобрить богов, и укрепляйте веру своей паствы с помощью проповедей.`, Price: 863},
	{ID: "1727127563", ImageURL: "https://shop.buka.ru/data/img_files/8975/screenshot/XhKcDnPNBf.jpg", Name: "Cyberpunk 2077", Description: `Cyberpunk 2077 — приключенческая ролевая игра с открытым миром, рассказывающая о киберпанке-наёмнике Ви и борьбе за жизнь в мегаполисе Найт-Сити. Мрачное будущее стало ещё более впечатляющим в улучшенной версии, в которую вошли новые дополнительные материалы. Создайте персонажа, выберите стиль игры и начните свой путь к высшей лиге, выполняя заказы, улучшая репутацию и оттачивая навыки. Ваши поступки влияют на происходящее и на весь город. В нём рождаются легенды. Какую сложат о вас?`, Price: 1513},
}

// обработчик для GET-запроса, возвращает список продуктов
func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для правильного формата JSON
	w.Header().Set("Content-Type", "application/json")
	// Преобразуем список заметок в JSON
	json.NewEncoder(w).Encode(products)
}

// обработчик для POST-запроса, добавляет продукт
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received new Product: %+v\n", newProduct)
	newProduct.ID = time.Now().String()
	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
}

//Добавление маршрута для получения одного продукта

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idStr := r.URL.Path[len("/Products/"):]

	// Ищем продукт с данным ID
	for _, Product := range products {
		if Product.ID == idStr {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Product)
			return
		}
	}

	// Если продукт не найден
	http.Error(w, "Product not found", http.StatusNotFound)
}

// удаление продукта по id
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	idStr := r.URL.Path[len("/Products/delete/"):]

	// Ищем и удаляем продукт с данным ID
	for i, Product := range products {
		if Product.ID == idStr {
			// Удаляем продукт из среза
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusNoContent) // Успешное удаление, нет содержимого
			return
		}
	}

	// Если продукт не найден
	http.Error(w, "Product not found", http.StatusNotFound)
}

// Обновление продукта по id
func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	idStr := r.URL.Path[len("/Products/update/"):]

	// Декодируем обновлённые данные продукта
	var updatedProduct Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ищем продукт для обновления
	for i, Product := range products {
		if Product.ID == idStr {

			products[i].ImageURL = updatedProduct.ImageURL
			products[i].Name = updatedProduct.Name
			products[i].Description = updatedProduct.Description
			products[i].Price = updatedProduct.Price

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products[i])
			return
		}
	}

	// Если продукт не найден
	http.Error(w, "Product not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/products", getProductsHandler)           // Получить все продукты
	http.HandleFunc("/products/create", createProductHandler)  // Создать продукт
	http.HandleFunc("/products/", getProductByIDHandler)       // Получить продукт по ID
	http.HandleFunc("/products/update/", updateProductHandler) // Обновить продукт
	http.HandleFunc("/products/delete/", deleteProductHandler) // Удалить продукт

	fmt.Println("Server is running on port 8080!")
	http.ListenAndServe(":8080", nil)
}
