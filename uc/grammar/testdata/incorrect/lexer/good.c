
/*							||
**	File for testing lexical analysis		||
**							||
**	This file would confuse a parser, but
        is 'lexically correct'.		                ||
*/

/* ** / ** */


// Simple tokens and single characters:

{ }       		                // until end-of-line comment
if else while                    	/* normal comment */
return && == != <= >=
char int void
			+ - * /		< > =		(,;) [was-colon] 

/* Comment with bad tokens: _ || | ++ # @ ...  */
// Ditto */ /* : _ || | ++ # @ ...  
// Identifiers and numbers:

17 -17 // No floats? -17.17e17 -17.17E-17  

	ponderosa Black Steel PUMPKIN AfterMath aBBaoN faT TRacKs

	K9 R23 B52 Track15 not4money 378 WHOIS666999SIOHM 
        was-floating-point-number

/* The following 'trap' should be correctly handled:

		* "2die4U" consists of the number '2' and the
		  identifier 'die4U'.
*/

        2die4U

//|| The following should all be regarded as identifiers:

	Function PrOceDuRE begIN eNd PrinT rEad iF THen StaTic
	ElSe wHilE Do reTurN noT AnD OR TrUE bOOl FalsE sizE


// It is legal to end the code like this, without an ending newline.
