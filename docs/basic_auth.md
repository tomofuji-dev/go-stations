# 基本認証テスト

## 正常系

```
curl -u admin:admin http://localhost:8080/healthz
```

- 200 OK

## 異常系

### 1. 間違ったユーザー ID とパスワードの場合

```
curl -u wronguser:wrongpass http://localhost:8080/healthz
```

#### 期待結果：

- 401 Unauthorized

### 2. 空のユーザー ID とパスワードの場合

```
curl -u : http://localhost:8080/healthz
```

#### 期待結果：

- 401 Unauthorized

### 3. ユーザー ID とパスワードを送信しない場合

```
curl http://localhost:8080/healthz
```

#### 期待結果：

- 401 Unauthorized
