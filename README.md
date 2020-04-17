### tl;dr
automated mouse bot to collect resource in Dofus

-----
### Usage

Run the bot simply with ```./dofus-bot```.

When a resource you want to collect is available, move your mouse over it and hit the key `+` on your keyboard. The key should be located at the right side of the top numeric bar.

For more mining effective result, please point at the top of the ore, like in the screen:

![Ore Position](https://github.com/AurelBec/dofus-bot/blob/master/doc/ore.png)

After this, the resource is registered, and every time the pixel will
change to the right color (i.e. the resource has been regenerated), the bot will click it!

-----
### Compilation and Run

```bash
make
./dofus-bot [-i] [-d] [-h]
```

-----
### Used libraries

- <a href="https://github.com/go-vgo/robotgo">RobotGo</a>