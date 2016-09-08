/* Initialisation* Includes the value for 2^0 = 1 and the ASCII code for Newline*/
>++++++++++>>+<+                                                    	
[
/* Write values on screen /*
   [+++++[>++++++++<-]
   >.
   <++++++[>--------<-]
   /*Check if there are less significant digits to display; If so display them using the same procedure* Repeat until A(0) is reached which contains 0*/
   +<<]
/* Write Newline */
   >.>
   /*Decrease the digit by one before starting the multiplication*/
   [->
   	/* Multiply A(3) by 2; storing the result in A(2) */
   	/* First round */
  	[<++>-
      	/*Second round*/
   [<++>-
         	/*Third round*/
[<++>-
   /*Fourth round*/
   [<++>-
  	/*Move 8 A(2) to 2 A(4) and 4(A3) and then*/
[<-------->>[-]++
/*Move 4 A(3) to 6 A(2) by decreasing A(3) before starting the duplication loop */
  	<-[<++>-]]
   ]
  	]
     	]
  	]<
   	/*Add A(2) to A(3); A(4) to A(5) etc */
  	[>+<-]
   	/*A(2)=1; A(4)=1; etc */
  	+>>]
   <<
]
