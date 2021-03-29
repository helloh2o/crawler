# /usr/bin

# Author ： xiaoyin
# Time : 2021/2/21 0021 22:26
# Describe ：
import requests
import execjs
from retrying import retry
import time


headers = {
    'user-agent':"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
}

with open('signature.js', 'r') as f:
    js_code = f.read()

def get_ttscid():
    '''
    获取到tt_scid
    :return:
    '''
    resp = requests.get('https://xxbg.snssdk.com/websdk/v1/getInfo?q=68VYzTEbWyXaOefcBmD3OK8A7LLiSdU1v5eqz%2B4A%2Far26FXYVkegrWdNvadH2SHGGHw4ttuacx3G%2BuUEuehWOy2%2FPTbwHTnXbW3YNbL2YF3mgKbEn2aPUmgfvGgvPqLxCDbgNNLLGN8uw%2FOYzIzvmy%2FNl98xolnEEAtA1wALZGMNEP2loBgJNG0%2BGaOZvBpi7pKJvS%2FYDuI3OL3dkLx%2B4Hl9cutIH8jivXsBYobTUrjItskM%2BAbZ9KnHTyblSBDnrN5i6kgzKoAC4n5x2SkvRcgoIvBgxZqSnG6M9iO%2FbcmN%2BOnFs0JzCWOQlPv9j%2BOeRmvwQtIMrBe5szUOEsJXZZOqEo%2Bru6%2BaiE5UKiUt15YaOxVu2%2FFqxbmX8S878jA0y4zz5Fto%2ByFjcprv8YSAwF9Qpo%2BcMsJnsJpN7h8bX01%2FtwXd45PlzjWtNeu5JsYQadH34RCyeBbFTjpVAVIhPwKzf10488zFCrU6X7evEq2cVS1eLDVJIck6UPYogsv2pClq7RhVZ5zb%2BZuMC5vb344a8FkdyWurnaQrNXzknsrn1ZXKvv98g633VUxi7Fon66%2FVpa9RoAeiNn2agYjyPkSoyE5daTOLEMFLTfxD0p5AShJoqiqZlzbPB6DmAeuvOwCDFUA2yij8SYhviHyjKc3kl08DsJ17TrNMScyWipj%2FRqSalCpNU5C%2Bnu7aA2CC3L8JmuwH8ZXqRh0kaF%2FkNRNvi1ef%2BhplzbPoubZFsay%2B5GXFileGYt9hQ%2F9H7rJEDEvsm2emYMvCoMKZbf1UTqRyKmiuf8q9PljPjfD74DXqF3XGlRlDCi%2BqCIWxEjy5kT6IWg6mTz9pBH4jGZvtNuKqbwMArULfE6YsHvq78PIZxxw4CDOYrIpRZ3mjVxnUljXHJoc8%2BdBN0x07LoOvXut3YXErAvlpYL94eG9u8XqurjpkQfGuAeDPV4rjmJAzy38QnO4vXk6hQslJTegC9Jxdox%2B1qSHumPU5e3%2Btkb6XwzFngm7yvquYm9iLyiPNEGIz2IybMc3ml22kM6jYoXVobzYO4cteBPDyBnACsVXBLDEbyGyjm3zUx1K1oEHEnkFG8TiYWII%2BJc7BCSoiMBPFoFKC6Vu0MnFoJtqIapbWo2z5%2FmoTbBWnBsePoHdxVU7F%2FNgx0r2N89Kx0r6jRNvi24X72Ngr8rRr%2FIho8oKn0ogJ&callback=_6361_1616147893310'.format(int(time.time()*1000)),headers = headers)
    return resp.cookies.get_dict()
#https://www.toutiao.com/api/pc/feed/?min_behot_time=0&category=news_finance&utm_source=toutiao&widen=1&tadrequire=true&_signature=_02B4Z6wo00f01LQoAiQAAIDBY-4vX4SD2ci0DQaAAE1jdWCYpIXjyskHsTfmP54ag7MBKGyiWK3s8wpCgPLWtoKb2A0OWz2vBI3vo26oA2ObZtu84N3rc.tlDhx9Da44WH9qiRVsWIny6zlgd4
#https://www.toutiao.com/api/pc/feed/?max_behot_time=1616135091&category=news_finance&utm_source=toutiao&widen=1&tadrequire=true&_signature=_02B4Z6wo00901pboUqwAAIDDQS5.1qHivoqWzVYAAMX6BqAMbpfBYmL2zpxm.P-LZiw84S5PXrm7KXmwlpuuJT3tGo9N9rbZKNd2Xv8eIIxyx0Wf24bMzMQVy5C-dCi6Bmsk1tk0QHQhr3YM37
@retry(stop_max_attempt_number=5, wait_fixed=5000)
def get_page(cookies,url,behot_time=0,key='min_behot_time'):
    '''
    获取列表页数据
    :return:
    '''
    scid = get_ttscid()
    tt_scid = 'tt_scid=' + scid['tt_scid']
    ctx = execjs.compile(js_code).call('get_page',url,tt_scid)
    params = {
         key: behot_time,
        'category': 'news_finance',
        'utm_source': 'toutiao',
        'widen': '1',
        'tadrequire': 'true',
        "_signature":ctx
    }
    cookies.update({'_signature':ctx,'tt_scid':scid['tt_scid'],'ttcid':scid['ttcid']})
    resp = requests.get(url=url,headers=headers,params=params,cookies = cookies)
    print(resp.url)
    result = resp.json()
    next_page = result['next']['max_behot_time']
    print(resp.json())
    time.sleep(5) # 休眠5秒
    get_page(cookies=cookies,url = url,behot_time=next_page,key='max_behot_time')
    return resp.json()

def main(url):
    resp = requests.get('https://www.toutiao.com/', headers=headers)
    
    cookies = resp.cookies.get_dict()
    get_page(cookies=cookies,url=url)

if __name__ == '__main__':
    main(url = 'https://www.toutiao.com/toutiao/api/pc/feed/')


