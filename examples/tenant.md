# Tenant

Tenant to projekt/zespol budujacy aplikacje skladajaca sie z komponentow - np. mikroserwisow, frontendu, bibliotek.
Tenant opisuje czlonkow zespolu, ich role (np. owner, approver, maintainer) oraz wspolne elementy dla wszystkich komponentow - projekt w Gitlab, grupy AD, Kubernetes namespace, itp.

Tenant jest tez wlascicielem wszystkich Resourcow wykorzystywanych przez indywidualne komponenty - np. bazy danych, topiki, kolejki.

W oryginalnym rozwiązaniu Tenant controller odpowiedzialny jest za stworzenie:

- System i Group w Backstage
- stworzenie grup w AD i dodanie czlonkow
- stworzenie wspolnych zasobow

Co należy wziąć pod uwagę:

- czy System zawsze powinien być tworzony - czy może Tenant jest częścią istniejącego Systemu
- Domain - czy tworzyć? Czy też część istniejącego
- jeśli tworzymy System/Domain - kwestia Groups i ich członkostwa. Czy Domain/System powinny mieć oddzielna grupę ownerową? Kto powinien być jej członkiem?
- kwestia ownershipu Componentow. Czy zespol rozwijajacy Component to osobny zespół?

Na przyklad, dla takiego Tenanta:

```yaml
apiVersion: idp.autopay.pl/v1alpha1
kind: Tenant
metadata:
  name: payments
  namespace: payments
spec:
  displayName: Payments
  description: |
    Payments services
  groups:
    - name: idp-payments-members
  members:
    - name: A B
      email: a.b@autopay.pl
      roles:
        - owner
        - maintainer
        - approver
    - name: C D
      email: c.d@autopay.pl
      roles:
        - maintainer
        - approver
    - name: F G
      email: f.g@autopay.pl
      roles:
        - developer
```

Stworzone zostana takie zasoby Backstage:

1. System

   ```yaml
   apiVersion: backstage.io/v1alpha1
   kind: System
   metadata:
     name: payments
     title: Payments
     description: Payments services
   spec:
     owner: group:default/payments-team
   ```

2. Glowna grupa zespolowa

   ```yaml
   apiVersion: backstage.io/v1alpha1
   kind: Group
   metadata:
     name: payments-team
     title: Payments team
   spec:
     type: team
     profile:
       email: payments@autopay.pl
       displayName: Payments
     children:
       - payments-team-owners
     members:
       - a-b
       - c-d
       - f-g
   ```

3. Podgrupy rolowe - zawierajace tylko czlonkow z dana rola
   ```yaml
   apiVersion: backstage.io/v1alpha1
   kind: Group
   metadata:
     name: payments-team-owners
     title: Payments team owners
   spec:
     type: owner-role-group
     members:
       - a-b
   ---
   apiVersion: backstage.io/v1alpha1
   kind: Group
   metadata:
     name: payments-team-approvers
     title: Payments team approvers
   spec:
     type: approver-role-group
     members:
       - a-b
       - c-d
   apiVersion: backstage.io/v1alpha1
   kind: Group
   metadata:
     name: payments-team-maintainers
     title: Payments team maintainers
   spec:
     type: maintainer-role-group
     members:
       - a-b
       - c-d
   ```
