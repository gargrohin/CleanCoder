#include<sys/types.h>
#include<unistd.h>
#include<stdio.h>

int main()
{
        printf("Think again");
        fflush(stdout);
            fork();
                printf("Hello World!\n");
                    return 0;


}
