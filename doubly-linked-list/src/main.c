#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct Node
{
    char *str;
    struct Node *next;
    struct Node *prev;
} Node;

typedef struct Dll
{
    Node *head;
    Node *tail;
} Dll;

Node *new_node(char *str)
{
    size_t len = strlen(str);
    char *s = malloc(len + 1);
    memcpy(s, str, len + 1);
    Node *n = malloc(sizeof *n);
    n->str = s;
    n->next = NULL;
    n->prev = NULL;
    return n;
}

void free_node(Node *n)
{
    free(n->str);
    free(n);
}

Dll *new_dll(void)
{
    Dll *d = malloc(sizeof *d);
    d->head = NULL;
    d->tail = NULL;
    return d;
}

void free_dll(Dll *d)
{
    for (Node *n = d->head; n;)
    {
        free(n->str);
        Node *next = n->next;
        free(n);
        n = next;
    }
    free(d);
    // free(d);
}

void insert_front(Dll *d, char *str)
{
    Node *n = new_node(str);
    if (d->head)
    {
        n->next = d->head;
        d->head->prev = n;
    }
    else
    {
        d->tail = n;
    }
    d->head = n;
}

void insert_back(Dll *d, char *str)
{
    Node *n = new_node(str);
    if (d->tail)
    {
        n->prev = d->tail;
        d->tail->next = n;
    }
    else
    {
        d->head = n;
    }
    d->tail = n;
}

void print_dll_forward(Dll *d)
{
    for (Node *n = d->head; n; n = n->next)
    {
        printf("%s ", n->str);
    }
    printf("\n");
}

void print_dll_backward(Dll *d)
{
    for (Node *n = d->tail; n; n = n->prev)
    {
        printf("%s ", n->str);
    }
    printf("\n");
}

Node *find(Dll *d, char *str)
{
    for (Node *n = d->head; n; n = n->next)
    {
        if (!strcmp(str, n->str))
        {
            return n;
        }
    }
    return NULL;
}

void delete(Dll *d, char *str)
{
    for (Node *n = d->head; n; n = n->next)
    {
        if (!strcmp(str, n->str))
        {
            if (n->prev && n->next) {
                n->prev = n->next;
                n->next->prev = n->prev;
            } else if (n->prev && !n->next) {
                d->tail = n->prev;
                n->prev->next = NULL;
            } else if(!n->prev && n->next) {
                d->head = n->next;
                n->next->prev = NULL;
            } else {
                d->head = NULL;
                d->tail = NULL;
            }
            free_node(n);
            return;
        }
    }
}

int main(void)
{
    Node *n = new_node("hooray");
    printf("%s\n", n->str);
    free_node(n);
    // printf("%s\n", n->str);

    Dll *d = new_dll();
    free_dll(d);

    d = new_dll();
    print_dll_forward(d);
    print_dll_backward(d);
    insert_front(d, "1");
    insert_front(d, "2");
    insert_front(d, "3");
    print_dll_forward(d);
    print_dll_backward(d);
    free_dll(d);

    d = new_dll();
    print_dll_forward(d);
    print_dll_backward(d);
    find(d, "asdf");
    insert_back(d, "5");
    insert_back(d, "6");
    insert_back(d, "7");
    insert_front(d, "2");
    print_dll_forward(d);
    print_dll_backward(d);
    Node *found = find(d, "6");
    if (found)
    {
        printf("found %s\n", found->str);
    }
    else
    {
        printf("not found\n");
    }
    found = find(d, "4");
    if (found)
    {
        printf("found %s\n", found->str);
    }
    else
    {
        printf("not found\n");
    }
    delete(d, "asdf");
    print_dll_forward(d);
    delete(d, "2");
    print_dll_forward(d);
    delete(d, "7");
    print_dll_forward(d);
    delete(d, "6");
    print_dll_forward(d);
    delete(d, "5");
    print_dll_forward(d);
    delete(d, "2");
    print_dll_forward(d);
    free_dll(d);

    return 0;
}