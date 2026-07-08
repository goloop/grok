# grok - довідник

Повний довідник пакета `grok`: клієнт, спільна модель `goloop/ai`, chat
completions (інтерфейс і нативний), стрімінг, генерація зображень і моделі.

Англійська версія: **[DOC.md](DOC.md)**.

## Зміст

- [Ментальна модель](#ментальна-модель)
- [Створення клієнта](#створення-клієнта)
- [Generate і Stream](#generate-і-stream)
- [Нативні chat completions](#нативні-chat-completions)
- [Інструменти, зображення й system-промпти](#інструменти-зображення-й-system-промпти)
- [Генерація зображень](#генерація-зображень)
- [Моделі](#моделі)
- [Опції та помилки](#опції-та-помилки)

## Ментальна модель

`grok.Client` реалізує `ai.Client` - провайдер-незалежний контракт із
`github.com/goloop/ai`. Спільні `Generate` і `Stream` покривають спільну основу
(чат із інструментами, зображеннями й стрімінгом), тож код проти інтерфейсу
працює з будь-яким провайдером.

Специфіка провайдера - у нативних методах: повний `ChatCompletion`, генерація
зображень, перелік моделей. Їх немає у спільному інтерфейсі. Формат обміну -
сумісний із chat completions.

```go
import (
	"github.com/goloop/ai"
	"github.com/goloop/grok"
)
```

## Створення клієнта

```go
c := grok.New(os.Getenv("XAI_API_KEY"))

c = grok.New(apiKey, grok.WithTimeout(30*time.Second))
```

Base URL за замовчуванням `https://api.x.ai/v1`. Наведіть `WithBaseURL` на
будь-який сумісний ендпоінт, щоб перевикористати клієнт.

## Generate і Stream

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    grok.ModelGrok4,
	System:   "You are concise.",
	Messages: []ai.Message{ai.UserText("Name three primary colors.")},
})
resp.Text()
resp.ToolCalls()
resp.Usage
```

`Stream` повертає `iter.Seq2[ai.Chunk, error]`: текстові дельти чанками з `Text`,
завершений виклик інструмента - чанком із `ToolCall`, фінальний чанк - `Done` і
`Usage`.

```go
for chunk, err := range c.Stream(ctx, req) {
	if err != nil {
		return err
	}
	fmt.Print(chunk.Text)
}
```

## Нативні chat completions

Для опцій, специфічних для провайдера, будуйте `ChatRequest` і викликайте
`ChatCompletion` чи `ChatCompletionStream`:

```go
resp, err := c.ChatCompletion(ctx, &grok.ChatRequest{
	Model:          grok.ModelGrok4,
	Messages:       []grok.ChatMessage{{Role: "user", Content: "as JSON"}},
	ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
})
```

Доступні `Tools`, `ToolChoice`, `Temperature`, `TopP`, `MaxTokens`, `Stop`, `N`,
`Seed`, `ResponseFormat`, `User`.

## Інструменти, зображення й system-промпти

Інструменти, зображення й system-промпти використовують спільні типи `ai`:
`ai.Tool`, `ai.Image`, `ai.ToolResult` і повідомлення `RoleSystem` або поле
`System`. Результати інструментів надсилаються назад повідомленнями `RoleTool`,
де `ai.ToolResult.ID` збігається з `ai.ToolUse.ID`. Вбудовані байти зображення
надсилаються як base64 data URI.

## Генерація зображень

```go
resp, err := c.GenerateImage(ctx, &grok.ImageRequest{
	Model: grok.ModelGrok2Image, Prompt: "a watercolor cat", N: 1,
})
resp.Data[0].URL // або B64JSON
```

## Моделі

```go
models, err := c.Models(ctx)
m, err := c.GetModel(ctx, grok.ModelGrok4)
```

## Опції та помилки

Опції: `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`,
`WithHeader`.

Невдала відповідь стає `*ai.APIError` зі `Status`, `Type`, `Code`, `Message` і
сирим тілом:

```go
var apiErr *ai.APIError
if errors.As(err, &apiErr) && apiErr.Status == http.StatusTooManyRequests {
	// backoff
}
```

Запити без моделі чи повідомлень падають до мережі з `ai.ErrNoModel` або
`ai.ErrNoMessages`.
