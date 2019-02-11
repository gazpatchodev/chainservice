OP_RETURN cache
---------------

When we query bitcoind, we do it by transaction hash / id.

So when someone asks for a transaction by itself, or for a specific vout on a transaction, we
will look on the disk for the transaction file and either return all the vouts, or only the
specified one.

If the transaction file does not exist on disk, we issue a `GetRawTransaction()` call to bitcoind
and process all the op return data.

When processing a raw transaction, we will store the raw JSON on disk in a file `<txid>.raw.json`
and the parsed results in `<txid>.opr.json`.  This way, if we wanted to re-parse the raw JSON with
a newer version of this tool, we could simply delete `*.opr.json` and walk the raw json files
without touching bitcoind.

> All JSON should be stored in multi-line format so tools like `grep` can be used.

The format of the `<txid>.opr.json` file will be:

```json
{
  "txid": "8bae12b5f4c088d940733dcd1455efc6a3a69cf9340e17a981286d3778615684",
  "blockHash": "000000000000000004c31376d7619bf0f0d65af6fb028d3b4a410ea39d22554c",
  "blockTime": "2014-06-30T05:45:09Z",
  "vout": [
    {
      "blocked": false // This field should not be included in any external responses.
      "n": 0,
      "value": 0,
      "hex": "6a13636861726c6579206c6f766573206865696469",
      "parts": [
        {
            "hex": "636861726c6579206c6f766573206865696469",
            "utf8": "charley loves heidi"
        }
      ]
    }
  ]   
}
```

Note the the vout section has an attribute "blocked".  This could be a solution to stop certain
vout content from being served by this service.

This is probably an independant mechanism specifically for this service.  Any downstream services,
like ElasticSearch, should have their own way of blocking specific content.


Logic flow
----------

1. Get TX for the request by looking for the file `<txid>.opr.json`

   1.1. If not found
   
      1.1.1. Get TX from bitcoind via `getRawTransaction()`

      1.1.2. Store the response in `<txid>.raw.json`

      1.1.3. Parse the vouts for OP_RETURN outputs
   
      1.1.4. Store the vouts in `<txid>.opr.json`
2. Return the JSON, filtering the vouts if necessary