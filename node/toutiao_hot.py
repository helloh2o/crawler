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
    resp = requests.get('https://xxbg.snssdk.com/websdk/v1/getInfo?q=3T0OfHdkf15dHX2XPmrz7pTbonbnvttvmtW2wY6hqtH4Plj0IM%2FSiKbwlTU8IlR6zq2S7sdE8C66WYVhj%2BfZ0XkbM0r527bx6hRe6BBX%2BwGPInWpIghqolf%2BnKAxjVaen0mx0%2Bh7fjrvX3iYYBNvAQuEH%2FXfqgM05n1k0S1JwUxNB5PRwbS6X3ksPKc2Xvayihl7esPgPbqYZbngkvoalTtHoixkAjbFMLABmPM5nvyewZYoVUUMm84smUKYwGJ8xG%2BelygbpUDMmTnXT8hCj1%2BbEgi24tPWrn%2Fothitd93oVFgFUONbyp6gxVS%2BmVL4K56vHhvS28Ov0AXP%2FAIcRArmgD07NX1pvQF4O2CFRdb4YrTimzUidRFnbgXgb7Siutl658dPuUGd87zjNyNBlV58VZknOyuLK%2BT6RwkyB%2FlBk8ltID%2BfLx3zeMUkBfXxfiIfmHxkBKejkK9vqZIYDLC2ljFpq3%2FEN1cREXaci5aAPxXoMCO7V1CrQkMRK2Ok9UdYjlscI3Z%2Bb419r8f3tUgrrHU6hlQ2VsDlVEB1xVj3tReUhOHNeJBkNV6zaJKJ9jh23G%2Bs8uHkr54xKsZQ%2FEsPVB6s5ZmkAZwmK13jy05x5xNTUTW%2BMXQyTAfCoRhexWeL9INqVRid5j3%2BpS%2FRIAxlIFOSeoOi31g4Ee6lLfLkzihItK%2B6P%2F52PCK7Rh8yYnA0qXF0OqS4fJtU2ctxdjOs8UPlwfjKW2WHZjT3jvS2epkQ8l9DGQlETOZzIkBXm6y%2FV%2BnVpliFyUTE4d3%2F8G786eBcU7bFSTM62QNFJMkbLOAbKsOmmMorP5X69HeU9wT62c3rPU4APReCOb68knzSqE3jrVOdnBQHU4FxL1Pt6vTz9cPyUzLZ8TWx3HFnjUHjK5ofPTtgjQrb0JClXen2ZFZpHyVPnBu%2Fv5y1TSl%2BqYU2DTtzeqE62K97TRgpJ5AhBJGjjhQY7d%2BSpKJBuoFVtYitJzrf6y5LV9aJygst35T0TpdntiCGMSaeKB6L6PFOSeHorUkUBwIpDefmhr2rBkiqdjLDUZkgxSuPXNbFMTha6qrFT%2FQf%2FelXbkpJttes2UJmM5MngCGl6uGF%2FZbcVQzld71%2FAsYQ%2BS7F5KECJ1Kk4CgP642N8NUN%2F4RG0489%2F3fx03ii0IKG030oRNvo23WV896V&callback=_9685_{}'.format(int(time.time()*1000)),headers = headers)
    return resp.cookies.get_dict()

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
        'category': '__all__',
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
    time.sleep(10) # 休眠5秒
    get_page(cookies=cookies,url = url,behot_time=next_page,key='max_behot_time')

    return resp.json()

def main(url):
    resp = requests.get('https://www.toutiao.com/ch/news_finance/', headers=headers)
    
    cookies = resp.cookies.get_dict()
    get_page(cookies=cookies,url=url)

if __name__ == '__main__':
    main(url = 'https://www.toutiao.com/toutiao/api/pc/feed/')


