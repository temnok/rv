The **APLIC (Advanced Platform-Level Interrupt Controller)** defined by the RISC-V Advanced Interrupt Architecture (AIA) is considerably more capable than the legacy PLIC. It supports both direct interrupt delivery and MSI generation for IMSICs. ([RISC-V Documentation][1])

If your goal is a **minimal implementation** (for example, a single hart with M-mode and optionally S-mode, and no IMSIC), you can implement only the **Direct Delivery** mode.

## Required functionality

For each interrupt source:

* Pending bit
* Enable bit
* Target hart
* Target privilege (M or S)
* Priority (or a fixed priority if you only support one level)

You also need:

* Domain enable (`domaincfg`)
* Interrupt source configuration (`sourcecfg`)
* Pending registers
* Enable registers
* Target registers

---

## Interrupt source state

For each source:

```c
struct irq {
    bool pending;
    bool enabled;
    bool delegated;      // M or S
    uint8_t priority;
    uint16_t hart;
};
```

For a minimal implementation, `priority` may be fixed to 1.

---

## Pending logic

An external device asserts an interrupt:

```c
irq[i].pending = true;
```

The APLIC searches for the highest-priority pending enabled interrupt for each hart.

If one exists:

```text
MEIP = 1
```

or

```text
SEIP = 1
```

depending on the target privilege.

---

## Claim/complete

Unlike the PLIC, the APLIC **does not have claim/complete registers** in Direct Delivery mode.

Instead:

1. Hardware raises `MEIP` or `SEIP`.
2. Software identifies the interrupt by reading the interrupt identity CSR (`mtopi` or `stopi` if implemented as part of AIA).
3. The interrupt source is cleared by writing to the device or by software clearing the pending bit in the APLIC, depending on the interrupt type. ([RISC-V Documentation][1])

---

## Registers worth implementing

A minimal set is:

| Register       | Purpose                     |
| -------------- | --------------------------- |
| `domaincfg`    | Enable the interrupt domain |
| `sourcecfg[i]` | Source type/configuration   |
| `setip[i]`     | Set pending                 |
| `clrip[i]`     | Clear pending               |
| `setie[i]`     | Enable interrupt            |
| `clrie[i]`     | Disable interrupt           |
| `target[i]`    | Hart and privilege target   |

Many optional registers (MSI configuration, IDCs, etc.) can be omitted if you only support Direct Delivery.

---

## Interrupt selection

Each cycle:

```c
for each source:
    if (pending &&
        enabled &&
        domain_enabled)
        candidate
```

Choose the highest-priority candidate for each hart.

Then:

```c
mip.MEIP = has_machine_interrupt;
mip.SEIP = has_supervisor_interrupt;
```

---

## Minimal datapath

```text
                 Device IRQs
                      |
                pending bits
                      |
                enable bits
                      |
               priority select
                      |
                target decode
                 /          \
             MEIP          SEIP
```

---

## What you can omit

For a very small implementation, you can omit:

* IMSIC support
* MSI generation
* Guest interrupts
* IDC registers
* Interrupt delegation between APLIC domains
* Complex source modes (edge/level variations beyond what your devices require)
* Full priority encoding (if a fixed priority is acceptable)

---

## Minimal implementation checklist

For a single-hart system:

* ✓ `domaincfg`
* ✓ One `sourcecfg` per interrupt
* ✓ Pending bits (`setip`/`clrip`)
* ✓ Enable bits (`setie`/`clrie`)
* ✓ One `target` register per interrupt
* ✓ Combinational selection of the highest-priority pending interrupt
* ✓ Drive `mip.MEIP` and/or `mip.SEIP`

This is sufficient to build a conforming Direct Delivery APLIC suitable for a simple single-hart system. MSI delivery, IMSIC integration, virtualization, and the richer AIA features can be added later without changing the basic pending/enable/target logic. ([RISC-V Documentation][1])

[1]: https://docs.riscv.org/reference/aia/v1.0/intro.html?utm_source=chatgpt.com "1.1. Introduction :: RISC-V Ratified Specifications Library"
