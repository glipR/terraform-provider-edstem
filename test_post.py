import requests
import os

res = requests.post("https://edstem.org/api/lessons/44601/slides", """-----------------------------264592028829639346041448524574
Content-Disposition: form-data; name="slide"

{"type":"document"}
-----------------------------264592028829639346041448524574--""".encode("utf-8"), headers={
    "Content-Type": "multipart/form-data; boundary=---------------------------264592028829639346041448524574",
    "X-Token": os.environ["EDSTEM_TOKEN"]
})

print(res.status_code)
print(res.content)
print(res)
