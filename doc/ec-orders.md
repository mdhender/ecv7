# Orders Reference
 
---
 
## 17.2.1. Combat Orders
 
```
scID "," "bombard"  "," scID "," pctAmt
scID "," "invade"   "," scID "," pctAmt
scID "," "raid"     "," scID "," pctAmt "," unitType
scID "," "support"  "," scID "," scID   "," pctAmt    -- attacking support
scID "," "support"  "," scID "," pctAmt               -- defensive support
```
 
---
 
## 17.2.2. Set Up Orders
 
```
"set up" "," locationID "," ("ship"|"colony") "," scID "," "transfer" ("," qty unitType)+ "," "END"
```
 
---
 
## 17.2.3. Assembly Orders
 
```
scID "," "assemble" "," qty "factories-" tlLevel "," unitType    -- factory assembly
scID "," "assemble" "," qty "mines-"     tlLevel "," locationID  -- mine assembly
scID "," "assemble" "," qty unitType                             -- other assembly
```
 
---
 
## 17.2.4. Dis-assembly Orders
 
Same structure as assembly orders, with `"dis-assemble"` replacing `"assemble"`.
 
```
scID "," "dis-assemble" "," qty "factories-" tlLevel "," unitType
scID "," "dis-assemble" "," qty "mines-"     tlLevel "," locationID
scID "," "dis-assemble" "," qty unitType
```
 
---
 
## 17.2.5. Build Change Orders
 
```
scID "," "build change" "," factoryGroupNo "," unitType
scID "," "build change" "," factoryGroupNo "," "retool"
scID "," "build change" "," "research"
```
 
---
 
## 17.2.6. Transfer Orders
 
```
scID "," "transfer" "," qty unitType "," scID
```
 
---
 
## 17.2.7. Mining Change Orders
 
```
scID "," "mining" "," miningGroupNo "," locationID
```
 
---
 
## 17.2.8. Market Orders
 
```
scID "," "sell" "," unitType    "," priceEach
scID "," "buy"  "," qty unitType "," priceEach
```
 
---
 
## 17.2.9. Survey Orders
 
```
scID "," "survey"
```
 
---
 
## 17.2.10. Probe Orders
 
```
scID "," "probe" "," orbitNo ("," orbitNo)*
```
 
---
 
## 17.2.11. Spy Orders
 
```
scID "," qty "," "check rebels"
scID "," qty "," "convert rebels"
scID "," qty "," "check for spies"
scID "," qty "," "attack spies"  "," scID
scID "," qty "," "incite rebels" "," scID
scID "," qty "," "information"   "," scID
```
 
---
 
## 17.2.12. News Release
 
```
"news" "," (planetID|tradeStationID) "," text ("," signature)?
```
 
---
 
## 17.2.13. Jump (Move) Orders
 
```
scID "," "move" "," locationID
```
 
`locationID` is either a planet number or a system coordinate (e.g. `4-6-19`).
 
---
 
## 17.2.14. Draft Orders
 
```
scID "," "draft"   "," qty unitType
scID "," "disband" "," qty
```
 
---
 
## 17.2.15. Pay Orders
 
```
scID "," "pay" "," wageAmt "," popClass
```
 
---
 
## 17.2.16. Ration Orders
 
```
scID "," "ration" "," pctAmt
```
 
---
 
## 17.2.17. Control Orders
 
```
playerID "," "control" "," systemID "," orbitNo
```
 
---
 
## 17.2.18. Un-Control Orders
 
```
playerID "," "un-control" "," systemID "," orbitNo
```
 
---
 
## 17.2.19. Naming Orders
 
```
playerID "," systemID "," planetNo "," name   -- planet (max 24 chars)
scID     "," name                             -- ship or colony (max 24 chars)
```
 
---
 
## 17.2.20. Trade Station Orders
 
```
stationID "," playerID "," ("permission granted"|"permission denied")
```
 
---
 
## 17.2.21. Colonizing Permission
 
```
playerID "," "permission to colonize" "," systemID "," planetNo
```
