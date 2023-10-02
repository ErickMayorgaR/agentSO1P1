//Info de los modulos
#include <linux/module.h>
//Info del kernel en tiempo real
#include <linux/kernel.h>
#include <linux/sched.h>
//Headers de los modulos
#include <linux/init.h>

// Header necesario para proc_fs
#include <linux/proc_fs.h>
// Para dar acceso al usuario
#include <asm/uaccess.h>
// Para manejar el directorio /proc
#include <linux/seq_file.h>
// Para get_mm_rss
#include <linux/mm.h>


MODULE_LICENSE("GPL");// Licencia del modulo
MODULE_DESCRIPTION("Modulo de CPU");
MODULE_DESCRIPTION("Módulo de Información de Memoria RAM");


static int escribir_archivo(struct seq_file *archivo, void *v) {
    struct sysinfo sys_info;
    si_meminfo(&sys_info);

    unsigned long total_memory = sys_info.totalram * sys_info.mem_unit;
    unsigned long free_memory = sys_info.freeram * sys_info.mem_unit;
    unsigned long used_memory = total_memory - free_memory;
    unsigned int usage_percent = (used_memory / total_memory) * 100;

    seq_printf(archivo, "Total_RAM: %lu\n", total_memory);
    seq_printf(archivo, "RAM_en_Uso: %lu\n", used_memory);
    seq_printf(archivo, "RAM_libre: %lu\n", free_memory);
    seq_printf(archivo, "Porcentaje_en_uso: %u\n", usage_percent);


    return 0;
}

//Funcion que se ejecutara cada vez que se lea el archivo con el comando CAT
static int al_abrir(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_archivo, NULL);
}

//Si el kernel es 5.6 o mayor se usa la estructura proc_ops
static struct proc_ops operaciones =
{
    .proc_open = al_abrir,
    .proc_read = seq_read
};

//Funcion a ejecuta al insertar el modulo en el kernel con insmod
static int _insert(void)
{
    proc_create("ram_201901758", 0, NULL, &operaciones);
    printk(KERN_INFO "201901758\n");
    return 0;
}

//Funcion a ejecuta al remover el modulo del kernel con rmmod
static void _remove(void)
{
    remove_proc_entry("ram_201901758", NULL);
    printk(KERN_INFO "Erick Ivan Mayorga Rodriguez\n");
}

module_init(_insert);
module_exit(_remove);
