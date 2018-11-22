%{
package selector
%}

%union {
    query Query
    path Path
    axis *Axis
    step *Step
    expr Expr
    args []Expr
    regexp *Regexp
    err error
    str string
    num int64
}

%type <query>   query
%type <path>    path
%type <axis>    axis
%type <step>    step
%type <expr>    filter func_call attr_ref query_param literal parenthesis
%type <expr>    expr expr1 expr2 expr3 expr4 expr5
%type <args>    func_args
%type <str>     axis_direction axis_type match logical_op bitwise_op relational_op arithmethical_op value

%token '[' ']' '(' ')' ':' '@' '.' ','

%token <regexp> REGEXP
%token <err>    ERR
%token <str>    ID STR TRUE FALSE NULL
%token <str>    '+' '-' '*' '/' '^' '%' '>' '<' '!' '~' '&' '|' '?'
%token <str>    LSHIFT RSHIFT AND OR EQ NE LTE GTE MATCH NONMATCH ELSE_OR
%token <num>    NUM

%start top

%%

top:
    query
    {
        querylex.(*queryLexerImpl).result = $1
    }
    ;

query:
    path
    {
        $$ = Query { $1 }
    }
|   query ',' path
    {
        $$ = append($1, $3)
    }
    ;

path:
    step
    {
        $$ = Path { $1 }
    }
|   path step
    {
        $$ = append($1, $2)
    }
    ;

step:
    match
    {
        $$ = &Step { Match: $1 }
    }
|   match filter
    {
        $$ = &Step { Match: $1, Filter: $2 }
    }
|   match '!' filter
    {
        $$ = &Step { Match: $1, Result: true, Filter: $3 }
    }
|   axis match
    {
        $$ = &Step { Axis: $1, Match: $2 }
    }
|   axis match filter
    {
        $$ = &Step { Axis: $1, Match: $2, Filter: $3 }
    }
|   axis match '!' filter
    {
        $$ = &Step { Axis: $1, Match: $2, Result: true, Filter: $4 }
    }
    ;

match:
    ID
|   STR
|   '*'
    ;

filter:
    '[' expr ']'
    {
        $$ = $2
    }
    ;

axis:
    axis_direction
    {
        $$ = &Axis { Dir: $1 }
    }
|   axis_direction axis_type
    {
        $$ = &Axis { $1, $2 }
    }
    ;

axis_direction:
    '/'             { $$ = "/" }
|   '/' '/'         { $$ = "//" }
|   '.' '/'         { $$ = "./" }
|   '.' '/' '/'     { $$ = ".//" }
|   '-' '/'         { $$ = "-/" }
|   '-' '/' '/'     { $$ = "-//" }
|   '+' '/'         { $$ = "+/" }
|   '+' '/' '/'     { $$ = "+//" }
|   '~' '/'         { $$ = "~/" }
|   '~' '/' '/'     { $$ = "~//" }
|   '.' '.' '/'     { $$ = "../" }
|   '.' '.' '/' '/' { $$ = "..//" }
|   '<' '/' '/'     { $$ = "<//" }
|   '>' '/' '/'     { $$ = ">//" }
    ;

axis_type:
    ':' ID  { $$ = $2 }
|   ':' STR { $$ = $2 }
    ;

expr:
    expr1
|   expr1 '?' expr1 ':' expr1
    {
        $$ = &Cond { $1, $3, $5 }
    }
|   expr1 ELSE_OR expr1
    {
        $$ = &Cond { $1, nil, $3 }
    }
    ;

expr1:
    expr2
|   expr2 logical_op expr2
    {
        $$ = &Binary { $1, $2, $3 }
    }
|   '!' expr2
    {
        $$ = &Unary { $1, $2 }
    }
    ;

logical_op:
    OR
|   AND
    ;

expr2:
    expr3
|   expr3 bitwise_op expr3
    {
        $$ = &Binary { $1, $2, $3 }
    }
|   '~' expr3
    {
        $$ = &Unary { $1, $2 }
    }
    ;

bitwise_op:
    '&'
|   '|'
|   LSHIFT
|   RSHIFT
    ;

expr3:
    expr4
|   expr4 relational_op expr4
    {
        $$ = &Binary{ $1, $2, $3 }
    }
|   expr4 MATCH REGEXP
    {
        $$ = &Binary{ $1, $2, $3 }
    }
|   expr4 NONMATCH REGEXP
    {
        $$ = &Binary{ $1, $2, $3 }
    }
    ;

relational_op:
    EQ
|   NE
|   LTE
|   GTE
|   '<'
|   '>'
    ;

expr4:
    expr5
|   expr5 arithmethical_op expr5
    {
        $$ = &Binary{ $1, $2, $3 }
    }
    ;

arithmethical_op:
    '+'
|   '-'
|   '*'
|   '/'
|   '%'
|   '^'
    ;

expr5:
    func_call
|   attr_ref
|   query_param
|   literal
|   parenthesis
    ;

func_call:
    ID '(' func_args ')'
    {
        $$ = &FuncCall { $1, $3 }
    }
    ;

func_args:
    /* empty */
    {
        $$ = nil
    }
|   expr
    {
        $$ = []Expr { $1 }
    }
|   func_args ',' expr
    {
        $$ = append($1, $3)
    }
    ;

attr_ref:
    '@' ID  { $$ = &WithAttr { $2 } }
|   '@' STR { $$ = &WithAttr { $2 } }
    ;

query_param:
    '{' ID '}'
    {
        $$ = QueryParam($2)
    }
    ;

literal:
    STR     { $$ = Str($1) }
|   NUM     { $$ = Num($1) }
|   REGEXP  { $$ = $1 }
|   value   { $$ = Keyword($1)}
    ;

value:
    TRUE
|   FALSE
|   NULL
    ;

parenthesis:
    '(' expr ')'
    {
        $$ = $2
    }
    ;

%%
