import requests

target_url = "https://club.jd.com/comment/productPageComments.action?productId=100008348530&score=0&sortType=5&page=0&pageSize=10"


def main():
    """
    main method, entry point
    :return: none
    """
    proxy = []
    with open("proxyPool.txt") as f:
        for p in f.readlines():
            proxy.append(p.strip())
    f.close()
    print(proxy)
if __name__ == '__main__':
    main()
