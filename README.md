# D-CMDB
####### from ops-updater
> modify heatbeat request response 5.20


###### 今天的心得：
1. cmdb 至少需要两组人实施
2. 一组人录入基础配置信息：位置，机房，购买时间等等
3. 一组人为应用维护人员：归类空闲主机与应用的关系
4. 如此才能合理的形成拓扑


### 模型0525
* agent 插入数据生成空闲
* 空闲机生成
* 手动关联业务
* 下线DOWN = 增加flag字段 = 状态变更，可以移动状态
* 告警关联
* 业务数
* 设备数
* 空闲设备
* 近一周新增设备情况
* 业务利用率
