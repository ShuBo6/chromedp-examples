package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"net/http"
	"time"
)

// var url = "https://shubo6.github.io/"
var url = "http://127.0.0.1:8765"

func main() {
	go testServer()
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body"),

		chromedp.ActionFunc(func(ctx context.Context) error {
			jsScroll := "window.scrollTo(0,document.documentElement.scrollHeight)"
			jsScroll1 := "window.scrollTo(0,0)"
			for i := 0; i < 10; i++ {
				chromedp.EvalAsValue(&runtime.EvaluateParams{
					Expression:    jsScroll,
					ReturnByValue: false,
				}).Do(ctx)
				time.Sleep(time.Second)
				chromedp.EvalAsValue(&runtime.EvaluateParams{
					Expression:    jsScroll1,
					ReturnByValue: false,
				}).Do(ctx)
			}

			return nil
		}),
		chromedp.OuterHTML("html", &res),
	)
	if err != nil {

		log.Panic(err)
	}
	fmt.Println(res)

}

func testServer() {
	s := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>无限加载示例</title>
    <style>
        #content {
            min-height: 2000px;
            padding: 20px;
            background-color: #f5f5f5;
            font-size: 24px;
            line-height: 1.5;
        }
    </style>
</head>
<body>

<div id="content">
    <p>向下滚动，当到达底部时，数字会增加。</p>
    <div id="number">0</div>
</div>

<script>
    // 监听滚动事件
    window.onscroll = function() {
		console.log("distanceFromBottom",distanceFromBottom)
		console.log("document.documentElement.scrollHeight",document.documentElement.scrollHeight)
		console.log("window.innerHeight",window.innerHeight)
		console.log("window.scrollY",window.scrollY)
        // 计算页面底部与当前视窗底部的距离
        var distanceFromBottom = document.documentElement.scrollHeight -
                                  window.innerHeight -
                                  window.scrollY;

        // 如果距离底部小于等于100px，增加数字
        if (distanceFromBottom <= 100) {
            var numberDiv = document.getElementById('number');
            var currentNumber = parseInt(numberDiv.textContent, 10);
            numberDiv.textContent = currentNumber + 1;
        }
    };
</script>

</body>
</html>`
	_ = http.ListenAndServe(":8765", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(s))
	}))

}
