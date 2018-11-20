%{
package selector

%}

%union {
    query Query
    path Path
    step *Step
    expr Expr
    args []Expr
    str string
    num int
    boolean bool
    dir AxisDirection
    axis *Axis
}

%type <query>   query
%type <path>    path
%type <step>    step
%type <expr>    filter expr condition logical bitwise relational arithmethical function_call attribute_ref query_parameter literal parenthesis
%type <args>    function_args
%type <str>     axis_type match logical_op bitwise_op relational_op arithmethical_op value
%type <dir>     axis_direction
%type <axis>    axis

%token '[', ']', '(', ')', ':', '@', '.', '"', '~', '=', ','

%token <str>    ID STR LAST POSITION '+', '-', '*', '/', '^', '%', '&', '|', '>', '<', '!', '~'
%token <num>    NUM
%token <boolean> BOOL

%right  '=' '!' '~' '?' ':'
%left   '&' '|' '>' '<'
%left   '+' '-'
%left   '*' '/' '%' '^'

%%

top:
    query

query:
    path
    {
        $$ = Query { $1 }
    }
|   query ',' path
    {
        $$ = append($1, $3)
    }

path:
    step
    {
        $$ = Path { $1 }
    }
|   path step
    {
        $$ = append($1, $2)
    }

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
        $$ = &Step { Match: $1, Filter: $3 }
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
        $$ = &Step { Axis: $1, Match: $2, Filter: $4 }
    }

match:
    ID
|   STR
|   '*'

filter:
    '[' expr ']'
    {
        $$ = $2
    }

axis:
    axis_direction
    {
        $$ = &Axis { Direction: $1 }
    }
|   axis_direction ':' axis_type
    {
        $$ = &Axis { $1, $3 }
    }

axis_direction:
    '/'             { $$ = DirectChild }
|   '/' '/'         { $$ = AnyDescendant }
|   '.' '/'         { $$ = CurrentDirectChild }
|   '.' '/' '/'     { $$ = CurrentAnyDescendant }
|   '-' '/'         { $$ = DirectLeftSibling }
|   '-' '/' '/'     { $$ = AnyLeftSibling }
|   '+' '/'         { $$ = DirectRightSibling }
|   '+' '/' '/'     { $$ = AnyRightSibling }
|   '~' '/'         { $$ = DirectLeftAndRightSibling }
|   '~' '/' '/'     { $$ = AnyLeftAndRightSibling }
|   '.' '.' '/'     { $$ = DirectParent }
|   '.' '.' '/' '/' { $$ = AnyParent }
|   '<' '/' '/'     { $$ = AnyPreceding }
|   '>' '/' '/'     { $$ = AnyFollowing }

axis_type:
    ID
|   STR

expr:
    condition
|   logical
|   bitwise
|   relational
|   arithmethical
|   function_call
|   attribute_ref
|   query_parameter
|   literal
|   parenthesis

condition:
    expr '?' expr ':' expr
    {
        $$ = &Condition { $1, $3, $5 }
    }
|   expr '?' ':' expr
    {
        $$ = &Condition { $1, nil, $4 }
    }

logical:
    expr logical_op expr
    {
        $$ = &Binary { $1, $2, $3 }
    }
|   '!' expr
    {
        $$ = &Unary { $1, $2 }
    }

logical_op:
    '&' '&' { $$ = "&&" }
|   '|' '|' { $$ = "||" }

bitwise:
    expr bitwise_op expr
    {
        $$ = &Binary { $1, $2, $3 }
    }
|   '~' expr
    {
        $$ = &Unary { $1, $2 }
    }

bitwise_op:
    '&'
|   '|'
|   '<' '<' { $$ = "<<" }
|   '>' '>' { $$ = ">>" }

relational:
    expr relational_op expr
    {
        $$ = &Binary{ $1, $2, $3 }
    }

relational_op:
    '=' '=' { $$ = "==" }
|   '!' '=' { $$ = "!=" }
|   '<' '=' { $$ = "<=" }
|   '>' '=' { $$ = ">=" }
|   '<'
|   '>'
|   '=' '~' { $$ = "=~" }
|   '!' '~' { $$ = "!~" }

arithmethical:
    expr arithmethical_op expr
    {
        $$ = &Binary{ $1, $2, $3 }
    }

arithmethical_op:
    '+'
|   '-'
|   '*'
|   '/'
|   '%'
|   '^'

function_call:
    ID '(' ')'
    {
        $$ = &FuncCall { $1, nil }
    }
|   ID '(' function_args ')'
    {
        $$ = &FuncCall { $1, $3 }
    }

function_args:
    expr
    {
        $$ = []Expr { $1 }
    }
|   function_args ',' expr
    {
        $$ = append($1, $3)
    }

attribute_ref:
    '@' ID  { $$ = &Attr { $2 } }
|   '@' STR { $$ = &Attr { $2 } }

query_parameter:
    '{' ID '}'
    {
        $$ = &QueryParam{ $2 }
    }

literal:
    STR     { $$ = Str($1) }
|   NUM     { $$ = Num($1) }
|   value   { $$ = Keyword($1)}

value:
    't' 'r' 'u' 'e'
    {
        $$ = "true"
    }
|   'f' 'a' 'l' 's' 'e'
    {
        $$ = "false"
    }
|   'n' 'u' 'l' 'l'
    {
        $$ = "null"
    }

parenthesis:
    '(' expr ')'
    {
        $$ = $2
    }

%%
