# libra һ����̬�ķ�����������

[English document](https://github.com/zhuCheer/libra/blob/master/README.md)�� [�����ĵ�](https://github.com/zhuCheer/libra/blob/master/README_CN.md)

## �������
- ��̬���ж��ַ�������
- ��̬����Դվ��ַ
- ֧�ֶ�̬�޸Ķ���Ӧͷ
- ���붼���ϸ�ĵ�Ԫ���ԣ��ɿ��ȸ�

�����ͨ���˰������ٹ���һ����̬�ĸ��ؾ����������Ŀǰ�����ָ��ؾ����㷨���ֱ����������ѯ����Ȩ��ѯ���˰�����ԭ�����빹���������������������������������ٵĻ��һ�����ؾ��������������ʹ�����ɡ�


## ���ٿ�ʼ


#### ���ذ�װ

�ڿ���̨����
`go get github.com/zhuCheer/libra`

#### ����ʾ��

���뵽���� example Ŀ¼��Ȼ������ example.go ���ɿ���һ��������������
```
> cd ../src/github.com/zhuCheer/libra/example
> go run example.go

```

��ʱ����ͨ����������� `http://127.0.0.1:5000` ���ɿ�������Ч����������ѯ�ķ�ʽ����`http://127.0.0.1:5001`, `http://127.0.0.1:5002` ������ http ����


#### ��ϸ����˵��
```
import "github.com/zhuCheer/libra"

    
// ע��һ�����������������������������ʵ� ip �Ͷ˿ڣ����ؾ������ͣ����Զ�����Ӧͷ ��������
// ���ؾ�������Ŀǰ�����ֿ�ѡ random:�����roundrobin:��ѯ��wroundrobin:��Ȩ��ѯ
srv := libra.NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)


// ��Ӵ���Ŀ��������� ip ��ַ�˿�
srv.GetBalancer().AddAddr("www.yourappdomain.com", "127.0.0.1:5001", 1)
srv.GetBalancer().AddAddr("www.yourappdomain.com", "127.0.0.1:5002", 1)


// ��������������
srv.Start()
```

#### ԭ�����

- �����ǵĸ���Ӧ�õ�������ֱ�ӽ������˴�������������ǿɳƴ˷���Ϊһ�����أ�
- ��������������������ͨ����ͬ���������ɷ��ʵ�������ӵ�ָ�� ip �ϣ������������ǻ��ж�����������������󽫰���ָ���ĸ��ؾ����㷨���е��ȣ�
- ���ʹ�������ͼ��ʾ��

![image](https://img.douyucdn.cn/data/yuba/weibo/2019/07/02/201907021730116899917826388.gif)

## ��ز�����������

```
import "github.com/zhuCheer/libra"
srv := libra.NewHttpProxySrv("127.0.0.1:5000", "roundrobin", nil)


// ������Ӧͷ
srv.ResetCustomHeader(map[string]string{"X-LIBRA": "the smart ReverseProxy"})

// �л����ؾ�������
srv.ChangeLoadType("random")


// ���Ŀ��������ڵ���Ϣ����Ӻ���Ч����������
srv.GetBalancer().AddAddr("www.yourappdomain.com","192.168.1.100:8081", 1)

// ɾ��Ŀ��������ڵ���Ϣ����Ӻ���Ч����������
srv.GetBalancer().DelAddr("www.yourappdomain.com","192.168.1.100:8081")

```