%{
package selector

%}

%union { 
    query Query
    str string
    num int
}

%type <query>   query query1 query2
%type <str>     binop

%token '/', '[', ']', '(', ')', '*', '=', ':', '@', '.', '"', '>', '<', '!'

%token <str>    '=', '>', '<', '!'
%token <str>    STR LAST POSITION 
%token <num>    NUM

%%

top:
    query    
    {

    }

query:
    query1
|   '/' query1
    {
        $$ = &DocElem { $2 }
    }
|   '/' '/' query1
    {
        $$ = &AllElems { $3 }
    }    
|   '.' '/' query1 
    {
        $$ = $3
    }
|   '.' '/' '/' query1 
    {
        $$ = &ChildElems { $4 }
    }


query1:
    query2
|   query2 '[' LAST '(' ')' ']'
    {
        $$ = &WithIndex { $1, -1 }
    }
|   query2 '[' POSITION '(' ')' binop NUM ']'
    {
        $$ = &WithPosition { $1, $6, $7 }
    }
|   query2 '[' NUM ']'
    {
        $$ = &WithIndex { $1, $3 }
    }
|   query2 '[' '@' STR ']'
    {
        $$ = &WithAttr { $1, $4 }
    }
|   query2 '[' '@' STR '=' '"' STR '"' ']'
    {
        $$ = &WithAttrValue { $1, $4, $7 }
    }

binop:
    '='
|   '!' '='
    {
        $$ = $1 + $2
    }
|   '>'
|   '<'
|   '>' '='
    {
        $$ = $1 + $2
    }
|   '<' '='
    {
        $$ = $1 + $2
    }

query2:
    STR 
    {
        $$ = &WithName { $1 }
    }
|   '(' query ')'
    {
        $$ = $2
    }

%%
