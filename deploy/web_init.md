1、 数据库的初始化
2、菜单管理添加

```
设置用户权限
菜单: view/user/users/users.vue 
api: 分组 users管理
POST   /users/createUsers       
DELETE /users/deleteUsers       
DELETE /users/deleteUsersByIds   
PUT    /users/updateUsers       
GET    /users/findUsers         
GET    /users/getUsersList       
GET    /users/getUsersPublic
```

```
设置客户管理权限
菜单: view/customers/customers/customers.vue
api: 分组 客户管理
POST   /customers/createCustomers
DELETE /customers/deleteCustomers
DELETE /customers/deleteCustomersByIds
PUT    /customers/updateCustomers
GET    /customers/findCustomers 
GET    /customers/getCustomersList
GET    /customers/getCustomersPublic
```