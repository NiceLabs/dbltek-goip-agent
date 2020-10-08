# DDNS Packet

Protocol: UDP

## Sent via GoIP

```plain
0000  86 70 14 13 01 38 4d 43  44 52 4d 31 38 30 34 37   .p...8MC DRM18047
0010  32 34 37                                           247

86 70 14 13 # magic header
01 # maybe sent flag
[Serial Number]
```

## Returns via "voipddns.com"

```plain
0000  86 70 14 13 00 32 31 38  2e 38 31 2e 33 37 2e 32   .p...218 .81.37.2
0010  30 33 09 77 77 77 2e 38  4d 43 44 52 4d 31 38 30   03.www.8 MCDRM180
0020  34 37 32 34 37 2e 63 6f  6d                        47247.co m

86 70 14 13 # magic header
00 # maybe returns flag
[IP Address]
09 # IP Address EOL flag
[Domain]
```
