package main

import(
	"os"
	"fmt"
	"log"
	"github.com/joho/godotenv"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// o'garuvchilar
var(
    bot 		*tgbotapi.BotAPI
    err 		error
    updChannel 	tgbotapi.UpdatesChannel
    update 		tgbotapi.Update
    updConfig 	tgbotapi.UpdateConfig
    botUser     tgbotapi.User
)

// menu tugmalar qatori
var (
	mainMenu = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🏠 Home"),
			tgbotapi.NewKeyboardButton("📹 Rec"),
            tgbotapi.NewKeyboardButton("🥸 Zohid"),
		 ),
	)
)

var (
    courseMenu = tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("Golang"),
            tgbotapi.NewKeyboardButton("Flutter"),
            tgbotapi.NewKeyboardButton("C++"),
         ),
         tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("C# (.Net)"),
            tgbotapi.NewKeyboardButton("Rust"),
            tgbotapi.NewKeyboardButton("Java"),
         ),
    )
)

type CourseSign struct {
    State int
    Name string
    Email string
    Telephone string
    Course string
}

var courseSignMap map[int]*CourseSign

func init() {
    courseSignMap = make(map[int]*CourseSign)
}

func main () {
    // load env
	if err = godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // connect tgbot
    bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_TOKEN"))
    if err != nil {
        log.Fatalf("bot init error: %v", err)
    }

    botUser, err := bot.GetMe()
    if err != nil {
        log.Fatalf("bot get me error: %v", err)
    }

    fmt.Printf("auth ok! bot is: %s\n", botUser.FirstName)

    // update configuring
    updConfig.Timeout = 60
    updConfig.Limit = 1
    updConfig.Offset = 0

    // channel for getting updates
    updChannel, err = bot.GetUpdatesChan(updConfig)
    if err != nil {
        log.Fatalf("update channel error: %v", err)
    }

   // button := tgbotapi.KeyboardButton(Text: "Golang")

    for {

        // get updates from channel
    	update = <- updChannel

        // check message
    	if update.Message != nil {
            // check message type
            if update.Message.IsCommand() {
                // get message text to variabel
                cmdText := update.Message.Command()
                if cmdText == "test" {
                    // create was sended message fields
                    msgConfig := tgbotapi.NewMessage(
                        update.Message.Chat.ID,
                        "test comand")
                    // send message    
                    bot.Send(msgConfig)  
                } else if cmdText == "menu" {
                   // keyMarkup := tgbotapi.NewInlineKeyboardMarkup()
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "main menu")
                    msg.ReplyMarkup = mainMenu
                    bot.Send(msg)
                } else {
                     msgConfig := tgbotapi.NewMessage(
                        update.Message.Chat.ID,
                        "/menu")
                    bot.Send(msgConfig)
                }  
            } else {
                // if message not command

                if update.Message.Text == mainMenu.Keyboard[0][1].Text {
                    
                    courseSignMap[update.Message.From.ID] = new(CourseSign)

                    fmt.Printf(
                        "message: %s\n",
                        update.Message.Text)
                    msgConfig := tgbotapi.NewMessage(
                        update.Message.Chat.ID,
                        "enter emil:")
                    msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
                    bot.Send(msgConfig)        
                } else {
                    cs, ok := courseSignMap[update.Message.From.ID]
                    if ok {
                        if cs.State == 0 {
                            cs.Email = update.Message.Text
                            cs.State = 1
                            msgConfig := tgbotapi.NewMessage(
                                update.Message.Chat.ID,
                                "create phone number:")
                            bot.Send(msgConfig)
                        } else if cs.State == 1 {
                            cs.Telephone = update.Message.Text
                            cs.State = 2
                            msgConfig := tgbotapi.NewMessage(
                                update.Message.Chat.ID,
                                "create course:")
                            msgConfig.ReplyMarkup = courseMenu
                            bot.Send(msgConfig)    
                        } else if cs.State == 2 {
                            cs.Course = update.Message.Text
                            cs.State = 2
                            msgConfig := tgbotapi.NewMessage(
                                update.Message.Chat.ID,
                                "OK! ")
                            msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
                            bot.Send(msgConfig)

                            delete(courseSignMap, update.Message.From.ID)
                            // post to site
                            fmt.Println(cs)
                            if err = SendPost(cs); err != nil {
                                fmt.Printf("send error: %v\n", err)
                            }

                        }
                        fmt.Printf("State: %+v\n", cs)
                    } else{
                        // other messages
                        msgConfig := tgbotapi.NewMessage(
                            update.Message.Chat.ID,
                            "OK! ")
                        msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
                        bot.Send(msgConfig)
                    }
                }
                /*fmt.Printf(
                    "message: %s\nfrom: %s\n",
                    update.Message.Text,
                    update.Message.From.FirstName,
                )

                msgConfig := tgbotapi.NewMessage(
                    update.Message.Chat.ID,
                    update.Message.From.FirstName + " " + update.Message.Text)
                bot.Send(msgConfig)*/                
            }
        } else {
            fmt.Printf("not message: %+v\n", update)
        }
    }

    bot.StopReceivingUpdates()
}