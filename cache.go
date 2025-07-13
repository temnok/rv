package rv

type Cache struct {
	VirtAddr, PhysAddr, Value int
}

func (c *Cache) Hit(a int) bool {
	return c.VirtAddr != 0 && c.VirtAddr == a
}

func (c *Cache) Clear() {
	c.VirtAddr = 0
}
