#include <stdio.h>
#include<sys/types.h>
#include <unistd.h>
int main(){
    pid_t pid1;
    pid1=0;
    printf("Parent\n");
    pid1=fork();
    if(pid1==0){
        printf("child\n");

    }
    else{
        printf("children\n");
    }
    return 0;
}
