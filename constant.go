package main

const helpmsg = `
Для того, что бы начать пользоваться ботом пришлите команду **/start**

Команды:

 **/newcat** - создает новую категорию расходов
*Пример:* 
/newcat Еда

**/newspnd** - создает новую трату. Слева сумма, справа комент к сумме. 
*Пример:* 
/newspnd 200,сникерс (запятая обязательно, даже если нет комента.) 


**/getcat** - вернет список твоих категорий 

**/reportmonth** - вернет картиночку с тратами за этот месяц.Если написать через пробел месяц и год, вернет за определенный месяц.(/reportmonth 01.2020) 

**/deletecat** - удалить категорию расходов 
`

const (
	startCommand       = "start"
	newCategory        = "newcat"
	deletecat          = "deletecategory"
	newSpendingCommand = "newspnd"
	report             = "reportMonth"
	helpCommand        = "help"
)
