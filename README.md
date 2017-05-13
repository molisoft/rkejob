### rkejob

纯粹是为了替代sidekiq而写的小开源项目。

因为最近 [锐壳云](https://www.rkecloud.com) 老是突然就cpu负载非常高，后来检查是内存不足，导致数据丢到swap时，系统异常慢，网站直接502~

#### 实现思路

因为不太想改动sidekiq中的业务，同时能把sidekiq关闭，使用这个程序替代。

所以决定之前的方式：
```
rails -> redis -> sidekiq(运行)
```

改为：

```
rails -> redis -> rkejob --post--> rails(运行)
```

相当于将任务又丢回了rails来运行。

所以这样的思路不适合用来跑运行时间非常长的任务，但是能降低非常多的内存，对于小内存跑rails时是很不错的！

如果业务量大了，将队列独立出来跑业务会是更好的选择。

#### 如何使用

rkejob：

程序同目录下创建config.yml文件（参考config.yml.example）

一定要配置

```
job:
    url: 您的网站url
```

rails：

sidekiq gem是要保留的，虽然已经不需要运行它了，但是需要依赖他来"推送队列"到redis中。

新增一个接受上面job url的controller action，处理rkejob推送过来的队列。

完成