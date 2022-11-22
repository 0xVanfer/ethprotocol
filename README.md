# Protocols

Not updated and not updating. 

Due to Coingecko request frequency limit, it will not always work properly.

Maybe will be deprecated.

```mermaid
graph LR
commonapys[apy info struct]
commonapys-->commonapy[apy]
commonapys-->commonapr[apr]
commonapys-->commonincy[incentive apy]
commonapys-->commonincr[incentive apr]

commonerc20s[erc20 info struct]
commonerc20s-->commonerc20address[address]
commonerc20s-->commonerc20symbol[symbol]
commonerc20s-->commonerc20decimals[decimals]

pro[protocol]

pro-->probasic[protocol basic info]
pro-->bps[lend pools]
pro-->lps[liquidity pools]
pro-->sps[stake pools]

probasic-->network
probasic-->proname[protocol name]
probasic-->client[eth client]
probasic-->geckokey[coingecko caller]

bps--aave-like-->avstoken[a,v,stoken]
bps--compound-like-->ctoken
bps-->bpstype[lend pool type]

bpstype-->aave-like
bpstype-->compound-like

avstoken-->erc20lendaave[basic erc20 info]-->erc201[erc20 info struct]
avstoken-->underlyinglendaave[underlying erc20 info]-->erc202[erc20 info struct]
avstoken-->apyslendaave[apy info]

apyslendaave-->apylendaave[apy info struct]

ctoken-->erc20lendcomp[basic erc20 info]-->erc203[erc20 info struct]
ctoken-->underlyinglendcomp[underlying erc20 info]-->erc204[erc20 info struct]
ctoken-->apyslendcompdepo[deposit apy info]-->apys1[apy info struct]
ctoken-->apyslendcomplend[borrow apy info]-->apys2[apy info struct]




```

## Targets

-   Pool

    -   tokens

    -   lp

    -   apys(apr)

    -   volume (day)

    -   tvl

    -   otherinfo(platypus)

    -   userinfo

        -   deposited

        -   ?

-   Lend

    -   atoken

        -   basic

        -   underlying

    -   vtoken

        -   basic

        -   underlying

    -   stoken

        -   basic

        -   underlying

    -   ctoken

        -   basic

        -   underlying

    -   deposit apys(apr)

    -   borrow apys(apr)

    -   collateral factor

    -   liquidation limit

    -   allow borrow

    -   allow collateral

    -   liquidation penalty

-   Stake

    -   tokens

    -   stake contract

    -   volume

    -   tvl

    -   apys(apr)
