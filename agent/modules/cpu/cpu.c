#define _GNU_SOURCE
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

#include <linux/sched/cpufreq.h>

// Estructura que almacena info del cpu
struct task_struct *cpu; 

// Almacena los procesos
struct list_head *lstProcess;
// Estructura que almacena info de los procesos hijos
struct task_struct *child;
unsigned long rss;


MODULE_LICENSE("GPL");// Licencia del modulo
MODULE_DESCRIPTION("Modulo de CPU");
MODULE_DESCRIPTION("Módulo de Información de Memoria RAM");


static int escribir_archivo(struct seq_file *archivo, void *v) {
    long long total_cpu_usage = 0;
    int max_cpu_usage = 0;

    bool is_first = true;

    for_each_process(cpu) {
    
        seq_printf(archivo, "{");
        
        seq_printf(archivo, "\"PID\": %d", cpu->pid);
        seq_printf(archivo, ",");
        seq_printf(archivo, "\"Nombre\": \"%s\"", cpu->comm);
        seq_printf(archivo, ",");
        seq_printf(archivo, "\"Status\": %lu", cpu->__state);
        seq_printf(archivo, ",");

         if (cpu->mm) {
            rss = get_mm_rss(cpu->mm) << PAGE_SHIFT;
            seq_printf(archivo, "\"Size\": %lu", rss);
        } else {
            seq_printf(archivo, "\"Size\": 0");

        }
        seq_printf(archivo, ",");

        seq_printf(archivo, "\"UID\": %d", cpu->cred->user->uid);
        seq_printf(archivo, "}\n");
    }   
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
    proc_create("cpu_201901758", 0, NULL, &operaciones);
    printk(KERN_INFO "201901758\n");
    return 0;
}

//Funcion a ejecuta al remover el modulo del kernel con rmmod
static void _remove(void)
{
    remove_proc_entry("cpu_201901758", NULL);
    printk(KERN_INFO "Erick Ivan Mayorga Rodriguez\n");
}

module_init(_insert);
module_exit(_remove);
