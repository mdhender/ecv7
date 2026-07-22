# Modern Orders Reference
 
This document defines the modernised order syntax for Epimethean Challenge.
It replaces the positional comma-separated format from the upstream document
with a keyword-tagged format that is easier for players to write and easier for
parsers to validate.
 
---
 
## Conventions
 
### General structure
 
Every order is one line. Fields are separated by spaces. String literals are
shown in `"double quotes"`. Optional fields are shown in `[square brackets]`.
Repeating fields are shown with `...`.
 
```
ENTITY scID  KEYWORD  keyword  value  keyword  value  ...
```
 
The leading `ENTITY scID` anchors every order to a specific ship or colony and
gives the parser a reliable sync point at the start of each line.
 
### IDs
 
| Token | Description |
|---|---|
| `scID` | Ship or colony ID (integer) |
| `playerID` | Player ID (integer) |
| `systemID` | System coordinate, e.g. `4-6-19` |
| `locationID` | Planet number (integer) or system coordinate |
| `groupNo` | Factory or mine group number (integer) |
 
### Values
 
| Token | Description |
|---|---|
| `qty` | Integer quantity, e.g. `1000` |
| `unitCode` | Unit code from the unit list, e.g. `SOL`, `FACT-3`, `MINE-2` |
| `tlLevel` | Tech level integer 0–10 |
| `pct` | Percentage as integer 0–100, e.g. `75` |
| `price` | Decimal price per unit, e.g. `0.25` |
| `wage` | Decimal wage multiplier, e.g. `1.2` |
| `name` | Quoted string, max 24 characters, e.g. `"Dragonfire"` |
| `text` | Quoted free text, e.g. `"Selling drives next turn"` |
 
### Percentage
 
All percentages are written as plain integers without a `%` symbol.
The keyword `commit` labels the value unambiguously.
 
---
 
## Orders
 
### 1. Combat Orders
 
One combat order per entity per turn.
 
```
ENTITY scID  BOMBARD   target scID  commit pct
ENTITY scID  INVADE    target scID  commit pct
ENTITY scID  RAID      target scID  commit pct  steal unitCode
ENTITY scID  SUPPORT   ally scID    commit pct
ENTITY scID  SUPPORT   ally scID    attacking scID  commit pct
```
 
Notes:
- `SUPPORT` with only `ally` and `commit` is defensive support.
- `SUPPORT` with `attacking scID` specifies which enemy the ally is attacking.
---
 
### 2. Set Up Orders
 
```
SETUP  location locationID  type ("SHIP"|"COPN"|"CENC"|"CORB")  source scID
  TRANSFER  qty unitCode
  TRANSFER  qty unitCode
  ...
END SETUP
```
 
Notes:
- `SETUP` / `END SETUP` are block delimiters — no `scID` prefix, as the new
  entity does not yet have one.
- Each `TRANSFER` line moves one unit type from the source entity.
- The block form makes multi-item set-ups easy to read and parse.
---
 
### 3. Assembly Orders
 
```
ENTITY scID  ASSEMBLE  qty unitCode
ENTITY scID  ASSEMBLE  qty FACT-tlLevel  produce unitCode
ENTITY scID  ASSEMBLE  qty MINE-tlLevel  deposit locationID
```
 
Notes:
- Plain `ASSEMBLE qty unitCode` covers all non-factory, non-mine assemblies.
- `produce` labels what the factory group will build.
- `deposit` labels the deposit the mine group will work.
---
 
### 4. Dis-assembly Orders
 
```
ENTITY scID  DISASSEMBLE  qty unitCode
ENTITY scID  DISASSEMBLE  qty FACT-tlLevel  produce unitCode
ENTITY scID  DISASSEMBLE  qty MINE-tlLevel  deposit locationID
```
 
Same structure as Assembly, keyword changed to `DISASSEMBLE`.
 
---
 
### 5. Build Change Orders
 
```
ENTITY scID  BUILDCHANGE  group groupNo  produce unitCode
ENTITY scID  BUILDCHANGE  group groupNo  retool
ENTITY scID  BUILDCHANGE  research
```
 
Notes:
- `produce unitCode` sets the new output of the factory group.
- `retool` begins the retooling process without specifying a new product yet.
- `research` redirects the entity's factories to research production.
---
 
### 6. Transfer Orders
 
```
ENTITY scID  TRANSFER  qty unitCode  to scID
```
 
---
 
### 7. Mining Change Orders
 
```
ENTITY scID  MININGCHANGE  group groupNo  deposit locationID
```
 
---
 
### 8. Market Orders
 
```
ENTITY scID  SELL  unitCode  price price
ENTITY scID  BUY   qty unitCode  price price
```
 
Notes:
- `SELL` lists a unit type at a price; quantity is not specified (sell all available at that price).
- `BUY` includes a quantity cap.
---
 
### 9. Survey Orders
 
```
ENTITY scID  SURVEY
```
 
---
 
### 10. Probe Orders
 
```
ENTITY scID  PROBE  orbit orbitNo  [orbit orbitNo  ...]
```
 
Multiple orbits may be listed, each prefixed with `orbit`.
 
---
 
### 11. Spy Orders
 
```
ENTITY scID  SPY  qty  CHECK REBELS
ENTITY scID  SPY  qty  CONVERT REBELS
ENTITY scID  SPY  qty  CHECK FOR SPIES
ENTITY scID  SPY  qty  ATTACK SPIES    from scID
ENTITY scID  SPY  qty  INCITE REBELS   at scID
ENTITY scID  SPY  qty  GATHER INFO     from scID
```
 
Notes:
- `qty` is the number of spy units committed.
- `from scID` / `at scID` identifies the foreign entity targeted.
---
 
### 12. News Release
 
```
NEWS  at locationID  text  [sig text]
```
 
Notes:
- `at locationID` is a planet or trade station ID.
- `text` is a quoted string containing the message body.
- `sig` is an optional quoted signature.
---
 
### 13. Move Orders
 
```
ENTITY scID  MOVE  to locationID
```
 
`locationID` is a planet number (integer) or a system coordinate (e.g. `4-6-19`).
 
---
 
### 14. Draft Orders
 
```
ENTITY scID  DRAFT    qty unitCode
ENTITY scID  DISBAND  qty unitCode
```
 
---
 
### 15. Pay Orders
 
```
ENTITY scID  PAY  wage wage  class ("USK"|"PRO"|"SOL")
```
 
Multiple `PAY` orders may be issued per entity, one per population class.
 
---
 
### 16. Ration Orders
 
```
ENTITY scID  RATION  pct pct
```
 
---
 
### 17. Control Orders
 
```
PLAYER playerID  CONTROL    system systemID  orbit orbitNo
PLAYER playerID  UNCONTROL  system systemID  orbit orbitNo
```
 
Notes:
- `PLAYER playerID` replaces `ENTITY scID` as the subject, since control is
  asserted by a player rather than a specific ship or colony.
---
 
### 18. Naming Orders
 
```
PLAYER playerID  NAMEP  system systemID  planet planetNo  name name
ENTITY scID      NAME   name name
```
 
Notes:
- `NAMEP` names a planet; `NAME` names a ship or colony.
- `name` is a quoted string, max 24 characters.
---
 
### 19. Trade Station Orders
 
```
ENTITY stationID  PERMIT  player playerID  ("GRANT"|"DENY")
```
 
---
 
### 20. Colonising Permission
 
```
PLAYER playerID  COLONIZE  system systemID  planet planetNo
```
 
---
 
## Summary Table
 
| Keyword | Subject | Required fields |
|---|---|---|
| `BOMBARD` | `ENTITY scID` | `target`, `commit` |
| `INVADE` | `ENTITY scID` | `target`, `commit` |
| `RAID` | `ENTITY scID` | `target`, `commit`, `steal` |
| `SUPPORT` | `ENTITY scID` | `ally`, `commit`, [`attacking`] |
| `SETUP` … `END SETUP` | *(block)* | `location`, `type`, `source`, `TRANSFER`s |
| `ASSEMBLE` | `ENTITY scID` | `qty unitCode`, [`produce`\|`deposit`] |
| `DISASSEMBLE` | `ENTITY scID` | `qty unitCode`, [`produce`\|`deposit`] |
| `BUILDCHANGE` | `ENTITY scID` | `group`, (`produce`\|`retool`\|`research`) |
| `TRANSFER` | `ENTITY scID` | `qty unitCode`, `to` |
| `MININGCHANGE` | `ENTITY scID` | `group`, `deposit` |
| `SELL` | `ENTITY scID` | `unitCode`, `price` |
| `BUY` | `ENTITY scID` | `qty unitCode`, `price` |
| `SURVEY` | `ENTITY scID` | *(none)* |
| `PROBE` | `ENTITY scID` | `orbit` × 1+ |
| `SPY` | `ENTITY scID` | `qty`, spy-op keyword, [`from`\|`at`] |
| `NEWS` | *(global)* | `at`, `text`, [`sig`] |
| `MOVE` | `ENTITY scID` | `to` |
| `DRAFT` | `ENTITY scID` | `qty unitCode` |
| `DISBAND` | `ENTITY scID` | `qty unitCode` |
| `PAY` | `ENTITY scID` | `wage`, `class` |
| `RATION` | `ENTITY scID` | `pct` |
| `CONTROL` | `PLAYER playerID` | `system`, `orbit` |
| `UNCONTROL` | `PLAYER playerID` | `system`, `orbit` |
| `NAMEP` | `PLAYER playerID` | `system`, `planet`, `name` |
| `NAME` | `ENTITY scID` | `name` |
| `PERMIT` | `ENTITY stationID` | `player`, (`GRANT`\|`DENY`) |
| `COLONIZE` | `PLAYER playerID` | `system`, `planet` |

