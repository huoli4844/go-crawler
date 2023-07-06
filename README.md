# go-crawler

这是一个真正的爬虫框架，它具有简单易用的编写方式，并提供出色的性能。灵感来源于scrapy，它融合了scrapy的优秀特点，并进行了进一步的优化。
这个框架的目标是提供一个功能强大、易于使用的爬虫框架，让您能够快速构建高效的网络爬虫应用。

[go-crawler](https://github.com/lizongying/go-crawler)
[document](https://pkg.go.dev/github.com/lizongying/go-crawler)

## Feature

* 编写爬虫变得简单而直观，同时保持强大的性能表现。
* 框架提供了简单且灵活的中间件机制，使您能够自定义和扩展爬虫的功能。
* 内置多种实用的中间件，让您能够更轻松地进行爬虫开发。这些中间件包括各种常用的功能，减少额外的工作量。
* 支持多种解析方式，您可以选择适合您的需求的解析器，从而简化页面解析的过程。
* 同样，框架也支持多种保存方式，让您根据具体情况选择最适合的保存方式，提高数据存储的灵活性。
* 提供了广泛的配置选项，让您能够灵活地调整爬虫的行为。
* 具备内置的devServer，轻松进行调试和开发。

## Usage

* 基本架构
    * Spider（爬虫）：Spider是发起请求并处理回调解析方法的核心组件。
      您可以为每个Spider设置一个名称，使用spider.SetName(name)方法进行命名。
    * BaseSpider（基础爬虫）：BaseSpider实现了Spider的公共方法，避免了在每个Spider中重复编写相同的代码。其中包括GetName和SetName等方法。
    * Crawler（爬虫处理器）：Crawler集成了Spider、Downloader（下载器）、Exporter（数据导出器）、Scheduler（调度器）等组件，是爬虫处理逻辑的中心。
      它负责协调这些组件的工作，并提供统一的接口供Spider调用。
    * 由于方法的继承关系，实际上Spider可以直接调用BaseSpider和Crawler中的部分方法，提高了代码的复用性和灵活性。
      这样的架构使得爬虫的编写更加简洁、可维护性更高，同时提供了基础功能和处理逻辑的封装，让开发者可以专注于具体的爬虫业务逻辑。
* crawler选项。
  这些选项提供了便捷的方式来配置爬虫框架的各个方面，包括模式、平台、浏览器、日志记录器、过滤器、下载器、导出器、中间件、Pipeline等，
  可以灵活地设置和定制爬虫的行为和功能。
    * WithMode 设置爬虫的模式（Mode），会执行SetMode
    * WithPlatforms 设置爬虫的平台（Platforms），会执行SetPlatforms
    * WithBrowsers 设置爬虫的浏览器（Browsers），会执行SetBrowsers
    * WithLogger 设置爬虫的日志记录器（Logger），会执行SetLogger
    * WithFilter 设置爬虫的过滤器（Filter），会执行SetFilter
    * WithDownloader 设置爬虫的下载器（Downloader），会执行SetDownloader
    * WithExporter 设置爬虫的导出器（Exporter），会执行SetExporter
    * WithMiddleware 设置爬虫的中间件（Middleware），会执行SetMiddleware
    * WithPipeline 设置爬虫的Pipeline，会执行SetPipeline
    * WithRetryMaxTimes 设置请求的最大重试次数（RetryMaxTimes），会执行SetRetryMaxTimes
    * WithTimeout 设置请求的超时时间（Timeout），会执行SetTimeout
    * WithInterval 设置请求的间隔时间（Interval），会执行SetInterval
    * WithOkHttpCodes 设置正常的HTTP状态码（OkHttpCodes），会执行SetOkHttpCodes
* Item类需要实现Item接口（可以组合ItemUnimplemented）
    * `GetReferer()` 可以获取到referer。
    * UniqueKey属性作为唯一键用于过滤和其他用途
    * Id属性用于保存主键
    * Data属性用于保存完整数据（必须是指针类型）
    * 内置Item实现：框架提供了一些内置的Item实现，如ItemNone、ItemCsv、ItemJsonl、ItemMongo、ItemMysql、ItemKafka等。
      您可以根据需要开启相应的Pipeline，以实现数据的保存功能。
* middleware/pipeline包括框架内置、公共自定义（internal/middlewares，internal/pipelines）和爬虫内自定义（和爬虫同module）。
* 对于中间件和Pipeline的顺序（order）是非常重要的。
  在框架中，确保不同中间件和Pipeline的order值不重复。如果有重复的order值，后面的中间件或Pipeline将替换前面的中间件或Pipeline。
  这种设计允许开发者对中间件和Pipeline的顺序进行精确控制，确保它们按照期望的顺序依次执行。
  通过合理设置order值，可以确保中间件和Pipeline按照特定的逻辑顺序进行处理，从而满足爬虫功能和业务需求。
* 在框架中，内置的中间件具有预定义的order值，这些order值是10的倍数。为了避免与内置中间件的order冲突，建议自定义中间件时选择不同的order值。
  内置中间件的order值：内置中间件的order值是10的倍数，例如10、20、30等。这些值是框架预留的，用于内置中间件的顺序控制。
  自定义中间件的order值：当您定义自己的中间件时，请选择避开内置中间件的order值。例如，您可以选择使用11、12、13等不同的order值来定义自定义中间件。
  中间件的顺序配置：根据中间件的功能和需求，按照预期的执行顺序进行配置。确保较低order值的中间件先执行，然后依次执行较高order值的中间件。
    * stats:10
        * 数据统计中间件，用于统计爬虫的请求、响应和处理情况。
        * 可以通过配置项enable_stats_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.StatsMiddleware), 10)`
    * dump:20
        * 控制台打印item.data中间件，用于打印请求和响应的详细信息。
        * 可以通过配置项enable_dump_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.DumpMiddleware), 20)`
    * proxy:30
        * 用于切换请求使用的代理。
        * 可以通过配置项enable_proxy_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.ProxyMiddleware), 30)`
    * robotsTxt:40
        * robots.txt支持中间件，用于支持爬取网站的robots.txt文件。
        * 可以通过配置项enable_robots_txt_middleware来启用或禁用，默认禁用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.RobotsTxtMiddleware), 40)`
    * filter:50
        * 过滤重复请求中间件，用于过滤重复的请求。默认只有在Item保存成功后才会进入去重队列。
        * 可以通过配置项enable_filter_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.FilterMiddleware), 50)`
    * file:60
        * 自动添加文件信息中间件，用于自动添加文件信息到请求中。
        * 可以通过配置项enable_file_middleware来启用或禁用，默认禁用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.FileMiddleware), 60)`
    * image:70
        * 自动添加图片的宽高等信息中间件
        * 用于自动添加图片信息到请求中。可以通过配置项enable_image_middleware来启用或禁用，默认禁用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.ImageMiddleware), 70)`
    * http:80
        * 创建请求中间件，用于创建HTTP请求。
        * 可以通过配置项enable_http_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.HttpMiddleware), 80)`
    * retry:90
        * 请求重试中间件，用于在请求失败时进行重试。
        * 默认最大重试次数为10。可以通过配置项enable_retry_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.RetryMiddleware), 90)`
    * url:100
        * 限制URL长度中间件，用于限制请求的URL长度。
        * 可以通过配置项enable_url_middleware和url_length_limit来启用和设置最长URL长度，默认启用和最长长度为2083。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.UrlMiddleware), 100)`
    * referer:110
        * 自动添加Referer中间件，用于自动添加Referer到请求中。
        * 可以根据referrer_policy配置项选择不同的Referer策略，DefaultReferrerPolicy会加入请求来源，NoReferrerPolicy不加入请求来源
        * 配置 enable_referer_middleware: true 是否开启自动添加referer，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.RefererMiddleware), 110)`
    * cookie:120
        * 自动添加Cookie中间件，用于自动添加之前请求返回的Cookie到后续请求中。
        * 可以通过配置项enable_cookie_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.CookieMiddleware), 120)`
    * redirect:130
        * 网址重定向中间件，用于处理网址重定向，默认支持301和302重定向。
        * 可以通过配置项enable_redirect_middleware和redirect_max_times来启用和设置最大重定向次数，默认启用和最大次数为1。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.RedirectMiddleware), 130)`
    * chrome:140
        * 模拟Chrome中间件，用于模拟Chrome浏览器。
        * 可以通过配置项enable_chrome_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.ChromeMiddleware), 140)`
    * httpAuth:150
        * HTTP认证中间件，通过提供用户名（username）和密码（password）进行HTTP认证。
        * 需要在具体的请求中设置用户名和密码。可以通过配置项enable_http_auth_middleware来启用或禁用，默认禁用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.HttpAuthMiddleware), 150)`
    * compress:160
        * 支持gzip/deflate解压缩中间件，用于处理响应的压缩编码。
        * 可以通过配置项enable_compress_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.CompressMiddleware), 160)`
    * decode:170
        * 中文解码中间件，支持对响应中的GBK、GB2312和Big5编码进行解码。
        * 可以通过配置项enable_decode_middleware来启用或禁用，默认启用。
        * 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.DecodeMiddleware), 170)`
    * device:180
        * 修改请求设备信息中间件，用于修改请求的设备信息，包括请求头（header）和TLS信息。目前只支持User-Agent随机切换。
        * 需要设置设备范围（Platforms）和浏览器范围（Browsers）。
        * Platforms: Windows/Mac/Android/Iphone/Ipad/Linux
        * Browsers: Chrome/Edge/Safari/FireFox
        * 可以通过配置项enable_device_middleware来启用或禁用，默认禁用。
* 启用方法：在NewApp中加入crawler选项`pkg.WithMiddleware(new(middlewares.DeviceMiddleware), 180)`
* 在爬虫框架中，Pipeline用于处理Item。
  通过配置不同的Pipeline，您可以方便地处理Item并将结果保存到不同的目标，如控制台、文件、数据库或消息队列中，以满足您的需求。
  以下是一些常用的Pipeline及其功能说明：
    * dump:10
        * 用于在控制台打印Item的详细信息。
        * 您可以通过配置项enable_dump_pipeline来控制是否启用该Pipeline，默认为启用。
    * file:20
        * 用于下载文件并保存到Item中。
        * 您可以通过配置项enable_file_pipeline来控制是否启用该Pipeline，默认为启用。
    * image:30
        * 用于下载图片并保存到Item中。
        * 您可以通过配置项enable_image_pipeline来控制是否启用该Pipeline，默认为启用。
    * filter:200
        * 用于对Item进行过滤。
        * 它可用于去重请求，需要在中间件同时启用filter。
        * 默认情况下，Item只有在成功保存后才会进入去重队列。
        * 您可以通过配置项enable_filter_pipeline来控制是否启用该Pipeline，默认为启用。
    * csv
        * 用于将结果保存到CSV文件中的Pipeline。
        * 需要在ItemCsv中设置`FileName`，指定保存的文件名称（不包含.csv扩展名）。
        * 您可以使用tag `column:""`来定义CSV文件的列名。
        * 启用方法是在创建应用程序时，将`pkg.WithPipeline(new(pipelines.CsvPipeline), 101)`添加到crawler选项中。
    * jsonLines
        * 用于将结果保存到JSON Lines文件中的Pipeline。
        * 需要在ItemJsonl中设置`FileName`，指定保存的文件名称（不包含.jsonl扩展名）。
        * 您可以使用tag `json:""`来定义JSON Lines文件的字段。
        * 启用方法是在创建应用程序时，将`pkg.WithPipeline(new(pipelines.JsonLinesPipeline), 102)`添加到crawler选项中。
    * mongo
        * 用于将结果保存到MongoDB中的Pipeline。
        * 需要在ItemMongo中设置`Collection`，指定保存的collection名称。
        * 您可以使用tag `bson:""`来定义MongoDB文档的字段。
        * 启用方法是在创建应用程序时，将`pkg.WithPipeline(new(pipelines.MongoPipeline), 103)`添加到crawler选项中。
    * mysql
        * 用于将结果保存到MySQL中的Pipeline。
        * 需要在ItemMysql中设置`Table`，指定保存的表名。
        * 您可以使用tag `column:""`来定义MySQL表的列名。
        * 启用方法是在创建应用程序时，将`pkg.WithPipeline(new(pipelines.MysqlPipeline), 104)`添加到crawler选项中。
    * kafka
        * 用于将结果保存到Kafka中的Pipeline。
        * 需要在ItemKafka中设置Topic，指定保存的主题名。
        * 您可以使用tag `json:""`来定义Kafka消息的字段。
        * 启用方法是在创建应用程序时，将`pkg.WithPipeline(new(pipelines.KafkaPipeline), 105)`添加到crawler选项中。
* 信号（Signal）是一种机制，用于在运行时处理外部发出的操作指令。通过捕获和处理信号，您可以实现对爬虫的控制和管理
* 在配置文件中配置全局的请求参数，并在具体的请求中可以覆盖这些全局配置，可以提供更灵活和细粒度的请求定制
* 框架内置了多个解析模块。这些解析模块提供了不同的选择器和语法，以适应不同的数据提取需求。您可以根据具体的爬虫任务和数据结构，选择适合您的解析模块和语法，从网页响应中准确地提取所需的数据。
    * query选择器 go-query是一个处理query选择器的库 [go-query](https://github.com/lizongying/go-query)
        * 通过调用`response.Query()`方法，您可以使用query选择器语法来从HTML或XML响应中提取数据。
    * xpath选择器 go-xpath是一个可用于XPath选择的库 [go-xpath](https://github.com/lizongying/go-xpath)
        * 通过调用`response.Xpath()`方法，您可以使用XPath表达式来从HTML或XML响应中提取数据。
    * gjson gjson是一个用于处理JSON的库
        * 通过调用`response.Json()`方法，您可以使用gjson语法从JSON响应中提取数据。
    * re选择器 go-re是一个处理正则的库 [go-re](https://github.com/lizongying/go-re)
        * 通过调用`response.Re()`方法，您可以使用正则表达式从响应中提取数据。
* 代理。它可以帮助爬虫在请求网站时隐藏真实IP地址。
    * 自行搭建隧道代理：您可以使用 [go-proxy](https://github.com/lizongying/go-proxy)
      等工具来搭建隧道代理。这些代理工具可以提供随机切换的代理功能，对调用方无感知，方便使用。
      您可以在爬虫框架中集成这些代理工具，以便在爬虫请求时自动切换代理。
      这是一个随机切换的隧道代理，调用方无感知，方便使用。后期会加入一些其他的调用方式，比如维持原来的代理地址。
    * 其他调用方式：除了随机切换的代理方式，后期可以考虑加入其他的调用方式。
      例如，保持原来的代理地址不变，或者使用其他代理池工具进行代理IP的管理和调度。这样可以提供更多灵活性和选择性，以满足不同的代理需求。
* 要提高爬虫的性能，您可以考虑关闭一些未使用的中间件或Pipeline，以减少不必要的处理和资源消耗。以下是一些建议：
    * 检查中间件：审查已配置的中间件，并根据需要禁用不使用的中间件。您可以在配置文件中进行修改，或者在爬虫的入口方法中进行相应的配置更改。
    * 禁用不需要的Pipeline：检查已配置的Pipeline，并禁用不需要的Pipeline。
      例如，如果您不需要保存结果到MongoDB，可以禁用MongoPipeline。
    * 评估性能影响：在禁用中间件或Pipeline之前，请评估其对爬虫性能的实际影响。确保禁用的部分不会对功能产生负面影响。
    * 可以禁用的配置:
        * enable_stats_middleware: false
        * enable_dump_middleware: false
        * enable_filter_middleware: false
        * enable_file_middleware: false
        * enable_image_middleware: false
        * enable_http_middleware: false
        * enable_retry_middleware: false
        * enable_referer_middleware: false
        * enable_http_auth_middleware: false
        * enable_cookie_middleware: false
        * enable_url_middleware: false
        * enable_compress_middleware: false
        * enable_decode_middleware: false
        * enable_redirect_middleware: false
        * enable_chrome_middleware: false
        * enable_device_middleware: false
        * enable_proxy_middleware: false
        * enable_robots_txt_middleware: false
        * enable_dump_pipeline: false
        * enable_file_pipeline: false
        * enable_image_pipeline: false
        * enable_filter_pipeline: false
* 文件下载
    * 如果您希望将文件保存到S3等对象存储中，需要进行相应的配置
    * Files下载
        * 在Item中设置Files请求：在Item中，您需要设置Files请求，即包含要下载的文件的请求列表。
          可以使用`item.SetFilesRequest([]*pkg.Request{...})`
          方法设置请求列表。
        * Item.Data结构：您的Item的Data字段需要实现pkg.File的切片，用于保存下载文件的结果。
          该字段的名称必须是Files，如`type DataFile struct {Files []*media.File}`。
    * Images下载
        * 在Item中设置Images请求：在Item中，您需要设置Images请求，即包含要下载的图片的请求列表。
          可以使用item.SetImagesRequest([]*pkg.Request{...})方法设置请求列表。
        * Item.Data结构：您的Item的Data字段需要实现pkg.Image的切片，用于保存下载图片的结果。
          该字段的名称必须是Images，如`type DataImage struct {Images []*media.Image}`。
* 爬虫结构
    * 建议按照每个网站（子网站）或者每个业务为一个spider。不必分的太细，也不必把所有的网站和业务都写在一个spider里
* 为了方便开发和调试，框架增加了本地devServer，并在`-m dev`模式下会默认启用。
  您可以自定义路由（routes），只需要实现`pkg.Route` 接口，并通过在Spider中调用`AddDevServerRoutes(...pkg.Route)`
  方法将其注册到devServer中。通过使用本地devServer，您可以在开发和调试过程中更方便地模拟和观察网络请求和响应，以及处理自定义路由逻辑。
  这为开发者提供了一个便捷的工具，有助于快速定位和解决问题。以下是devServer的一些特性：
    * 支持http和https，您可以通过设置`dev_server`选项来指定devServer的URL。
      `http://localhost:8081`表示使用HTTP协议，`https://localhost:8081`表示使用HTTPS协议。
    * 默认显示JA3指纹。JA3是一种用于TLS客户端指纹识别的算法，它可以显示与服务器建立连接时客户端使用的TLS版本和加密套件等信息。
    * 您可以使用tls工具来生成服务器的私钥和证书，以便在devServer中使用HTTPS。tls工具可以帮助您生成自签名的证书，用于本地开发和测试环境。
    * devServer内置了多种handler，这些handler提供了丰富的功能，可以模拟各种网络情景，帮助进行开发和调试。
      您可以根据需要选择合适的handler，并将其配置到devServer中，以模拟特定的网络响应和行为。以下是可用的handler列表及其功能：
        * BadGatewayHandler 模拟返回502状态码
        * Big5Handler 模拟使用big5编码
        * CookieHandler 模拟返回cookie
        * DeflateHandler 模拟使用Deflate压缩
        * FileHandler 模拟输出文件
        * Gb2312Handler 模拟使用gb2312编码
        * Gb18030Handler 模拟使用gb18030编码
        * GbkHandler 模拟使用gbk编码
        * GzipHandler 模拟使用gzip压缩
        * HelloHandler 打印请求的header和body信息
        * HttpAuthHandler 模拟http-auth认证
        * InternalServerErrorHandler 模拟返回500状态码
        * OkHandler 模拟正常输出，返回200状态码
        * RateLimiterHandler 模拟速率限制，目前基于全部请求，不区分用户。可与HttpAuthHandler配合使用。
        * RedirectHandler 模拟302临时跳转，需要同时启用OkHandler
        * RobotsTxtHandler 返回robots.txt文件

### args

通过配置环境变量和启动参数，您可以更灵活地配置和控制爬虫的行为，包括选择配置文件、指定入口方法、传递额外参数以及设定启动模式。这样的设计可以提高爬虫的可配置性和可扩展性，使得爬虫框架更适应各种不同的应用场景。

* CRAWLER_CONFIG_FILE -c 配置文件路径，必须进行配置。该配置文件包含了爬虫的各项配置信息。
* CRAWLER_START_FUNC -f 入口方法名称，默认为Test。您可以根据实际需要自定义入口方法，用于启动爬虫的执行流程。
* CRAWLER_ARGS -a 额外的参数，以JSON字符串的形式提供。这些参数可以在入口方法调用时使用，用于进一步定制爬虫的行为。该参数是非必须项，根据具体需求进行配置。
* CRAWLER_MODE -m 启动模式，默认为test。您可以根据需要配置不同的模式，如dev、prod等，以适应不同的开发和生产环境。

### config

数据库相关配置：

* mongo_enable: 是否启用MongoDB。
* mongo.example.uri: MongoDB的连接URI。
* mongo.example.database: MongoDB的数据库名称。
* mysql_enable: 是否启用MySQL。
* mysql.example.uri: MySQL的连接URI。
* mysql.example.database: MySQL的数据库名称。
* redis_enable: 是否启用Redis。
* redis.example.addr: Redis的地址。
* redis.example.password: Redis的密码。
* redis.example.db: Redis的数据库编号。
* s3_enable: 是否启用S3对象存储（如COS、OSS、MinIO等）
* s3.example.endpoint:  S3的终端节点
* s3.example.region: S3的区域。
* s3.example.id: S3的身份标识。
* s3.example.key: S3的身份密钥。
* s3.example.bucket: S3的存储桶名称。
* kafka_enable: 是否启用Kafka。
* kafka.example.uri: Kafka的连接URI。

日志和日志文件配置：

* log.filename: 日志文件路径。可以使用"{name}"替换为-ldflags。
* log.long_file: 如果设置为true，则记录完整文件路径。
* log.level: 日志级别，可选DEBUG/INFO/WARN/ERROR。

其他配置项：

* proxy.example: 代理配置。
* request.concurrency: 请求并发数。
* request.interval: 请求间隔时间（毫秒）。如果设置为0，则使用默认间隔时间（1000毫秒）。如果设置为负数，则为0。
* request.timeout: 请求超时时间（秒）。
* request.ok_http_codes: 请求正常的HTTP状态码。
* request.retry_max_times: 请求重试的最大次数，默认为10。
* request.http_proto: 请求的HTTP协议。
* dev_server: 开发服务器（devServer）的地址。
* enable_ja3: 是否显示devServer的JA3指纹。

其他中间件和pipeline相关配置。

* enable_stats_middleware: 是否开启统计中间件，默认为true。
* enable_dump_middleware: 是否开启打印请求/响应中间件，默认为true。
* enable_filter_middleware: 是否开启过滤中间件，默认为true。
* enable_file_middleware: 是否开启文件处理中间件，默认为true。
* enable_image_middleware: 是否开启图片处理中间件，默认为true。
* enable_http_middleware: 是否开启HTTP请求中间件，默认为true。
* enable_retry_middleware: 是否开启请求重试中间件，默认为true。
* enable_referer_middleware: 是否开启Referer中间件，默认为true。
* referrer_policy: 设置来源政策，可选值为DefaultReferrerPolicy（默认值）和NoReferrerPolicy。
* enable_http_auth_middleware: 是否开启HTTP认证中间件，默认为false。
* enable_cookie_middleware:  是否开启Cookie中间件，默认为true。
* enable_url_middleware: 是否开启URL长度限制中间件，默认为true。
* url_length_limit: URL的最大长度限制，默认为2083。
* enable_compress_middleware: 是否开启响应解压缩中间件（支持gzip/deflate），默认为true。
* enable_decode_middleware: 是否开启中文解码中间件（支持GBK、GB2312、Big5编码），默认为true。
* enable_redirect_middleware: 是否开启重定向中间件，默认为true。
* redirect_max_times: 重定向的最大次数，默认为1。
* enable_chrome_middleware: 是否开启Chrome模拟中间件，默认为true。
* enable_device_middleware: 是否开启设备模拟中间件，默认为false。
* enable_proxy_middleware: 是否开启代理中间件，默认为true。
* enable_robots_txt_middleware: 是否开启robots.txt支持中间件，默认为false。
* enable_dump_pipeline: 是否开启打印Item Pipeline，默认为true。
* enable_file_pipeline: 是否开启文件下载Pipeline，默认为true。
* enable_image_pipeline: 是否开启图片下载Pipeline，默认为true。
* enable_filter_pipeline: 是否开启过滤Pipeline，默认为true。
* scheduler: 调度方式，默认为memory（单机调度），可选值为memory、redis、kafka。选择redis或kafka后可以实现集群调度。
* filter: 过滤方式，默认为memory（内存过滤），可选值为memory、redis。

## Example

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lizongying/go-crawler/pkg"
	"github.com/lizongying/go-crawler/pkg/app"
	"github.com/lizongying/go-crawler/pkg/devServer"
)

type ExtraOk struct {
	Count int
}

type DataOk struct {
	Count int
}

type Spider struct {
	pkg.Spider
	logger pkg.Logger
}

func (s *Spider) ParseOk(ctx context.Context, response *pkg.Response) (err error) {
	var extra ExtraOk
	_ = response.Request.GetExtra(&extra)

	item := pkg.ItemNone{
		ItemUnimplemented: pkg.ItemUnimplemented{
			Data: &DataOk{
				Count: extra.Count,
			},
		},
	}
	err = s.YieldItem(ctx, &item)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	if extra.Count > 0 {
		return
	}

	requestNext := new(pkg.Request)
	requestNext.Url = response.Request.Url
	requestNext.Extra = &ExtraOk{
		Count: extra.Count + 1,
	}
	requestNext.CallBack = s.ParseOk
	err = s.YieldRequest(ctx, requestNext)
	if err != nil {
		s.logger.Error(err)
	}
	return
}

func (s *Spider) TestOk(ctx context.Context, _ string) (err error) {
	// mock server
	s.AddDevServerRoutes(devServer.NewOkHandler(s.logger))

	request := new(pkg.Request)
	request.Url = fmt.Sprintf("%s%s", s.GetDevServerHost(), devServer.UrlOk)
	request.Extra = &ExtraOk{}
	request.CallBack = s.ParseOk
	err = s.YieldRequest(ctx, request)
	if err != nil {
		s.logger.Error(err)
	}
	return
}

func NewSpider(baseSpider pkg.Spider) (spider pkg.Spider, err error) {
	if baseSpider == nil {
		err = errors.New("nil baseSpider")
		return
	}

	spider = &Spider{
		Spider: baseSpider,
		logger: baseSpider.GetLogger(),
	}
	spider.SetName("test-ok")

	return
}

func main() {
	app.NewApp(NewSpider).Run()
}

```

### Test

```shell
go run cmd/testOkSpider/*.go -c example.yml -f TestOk -m dev

```

更多示例可以按照以下项目

[go-crawler-example](https://github.com/lizongying/go-crawler-example)

```shell
git clone github.com/lizongying/go-crawler-example
```

## TODO

* middlewares
    * downloadtimeout

* AutoThrottle
* cron
* max request limit?
* multi-spider
* devServer独立拆分

