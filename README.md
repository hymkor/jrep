JREP - Format-Preserving JSON Extractor (JSON生テキスト抽出ツール)
=======================================

```
> jrep .[1628].entities.user_mentions ../jegan/__twitter/tweet.js

    "user_mentions" : [ {
      "name" : "職業亡者",
      "screen_name" : "hymkor",
      "indices" : [ "3.0", "10.0" ],
      "id_str" : "94929449",
      "id" : "9.4929449E7"
    } ],

> jrep .[1628].entities.user_mentions[0] ../jegan/__twitter/tweet.js
 {
      "name" : "職業亡者",
      "screen_name" : "hymkor",
      "indices" : [ "3.0", "10.0" ],
      "id_str" : "94929449",
      "id" : "9.4929449E7"
    }

> jrep .[1628].entities.user_mentions[0].indices ../jegan/__twitter/tweet.j

      "indices" : [ "3.0", "10.0" ],

> jrep .[1628].entities.user_mentions[0].indices[0] ../jegan/__twitter/tweet.js
 "3.0",
```
