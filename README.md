# StarRocks

## Version

### v0.0.1
* 增加了解析StarRocks审计日志中SQL的功能，用于获取表访问记录。

### v0.0.2
* 完成StarRocks中所有和table直接相关非授权语句的解析.
* 不支持间接语句的解析，例如：cancel load 语句.
* 不支持like语句解析，例如：show tables like "%tableName%".
